package log

import (
	"sync"
	"time"

	"github.com/efritz/glock"
)

// FieldRollup is a field assigned to the last message in a
// window. Its value is equal to the number of messages in the
// window before it was flushed.
const FieldRollup = "rollup-multiplicity"

type (
	// RollupAdapter provides a way to throttle equivalent messages. Messages begin
	// a roll-up when a second messages with an identical format string is seen in
	// the same window period. All remaining messages logged within that period are
	// captured and emitted as a single message at the end of the window period. The
	// fields and args are equal to the first rolled-up message.
	RollupAdapter struct {
		logger         Logger
		clock          glock.Clock
		windowDuration time.Duration
		windows        map[string]*logWindow
		mutex          sync.RWMutex
	}

	logWindow struct {
		stashed *logMessage
		start   time.Time
		count   int
		mutex   sync.RWMutex
	}
)

//
// Shim

var _ logShim = &RollupAdapter{}

// NewRollupAdapter creates a RollupAdapter which wraps the given logger.
func NewRollupAdapter(logger Logger, windowDuration time.Duration) *RollupAdapter {
	return newRollupAdapter(logger, glock.NewRealClock(), windowDuration)
}

func newRollupAdapter(logger Logger, clock glock.Clock, windowDuration time.Duration) *RollupAdapter {
	return &RollupAdapter{
		logger:         logger,
		clock:          clock,
		windowDuration: windowDuration,
		windows:        map[string]*logWindow{},
	}
}

func (a *RollupAdapter) WithFields(fields Fields) logShim {
	if len(fields) == 0 {
		return a
	}

	return newRollupAdapter(
		a.logger.WithFields(fields),
		a.clock,
		a.windowDuration,
	)
}

func (a *RollupAdapter) Log(level LogLevel, format string, args ...interface{}) {
	a.LogWithFields(level, nil, format, args)
}

func (a *RollupAdapter) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	if a.getWindow(format).record(a.logger, a.clock, a.windowDuration, level, fields, format, args...) {
		// Not rolling up, log immediately
		logWithFields(a.logger, level, fields, format, args...)
	}
}

func (a *RollupAdapter) getWindow(format string) *logWindow {
	a.mutex.RLock()
	if window, ok := a.windows[format]; ok {
		a.mutex.RUnlock()
		return window
	}

	a.mutex.RUnlock()
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if window, ok := a.windows[format]; ok {
		return window
	}

	window := &logWindow{}
	a.windows[format] = window
	return window
}

func (a *RollupAdapter) Sync() error {
	return a.logger.Sync()
}

//
// Log Window

func (w *logWindow) record(logger Logger, clock glock.Clock, windowDuration time.Duration, level LogLevel, fields Fields, format string, args ...interface{}) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	now := clock.Now()

	if remaining := windowDuration - now.Sub(w.start); w.start != (time.Time{}) && remaining > 0 {
		w.count++

		if w.count == 1 {
			ch := clock.After(remaining)

			go func() {
				<-ch

				w.mutex.Lock()
				w.flush(logger)
				w.mutex.Unlock()
			}()
		}

		return false
	}

	w.flush(logger)

	w.count = 0
	w.start = now
	w.stashed = &logMessage{
		level:  level,
		fields: fields,
		format: format,
		args:   args,
	}

	return true
}

func (w *logWindow) flush(logger Logger) {
	if w.stashed == nil || w.count <= 1 {
		return
	}

	fields := w.stashed.fields
	if fields == nil {
		fields = Fields{}
	}

	// Set replay field on message
	fields[FieldRollup] = w.count

	logWithFields(
		logger,
		w.stashed.level,
		fields,
		w.stashed.format,
		w.stashed.args...,
	)

	w.stashed = nil
}
