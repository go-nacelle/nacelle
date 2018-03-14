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
	rollupShim struct {
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

var _ logShim = &rollupShim{}

// NewRollupAdapter returns a logger with functionality to throttle similar messages. Messages begin
// a roll-up when a second messages with an identical format string is seen in
// the same window period. All remaining messages logged within that period are
// captured and emitted as a single message at the end of the window period. The
// fields and args are equal to the first rolled-up message.
func NewRollupAdapter(logger Logger, windowDuration time.Duration) Logger {
	return adaptShim(newRollupShim(logger, glock.NewRealClock(), windowDuration))
}

func newRollupShim(logger Logger, clock glock.Clock, windowDuration time.Duration) *rollupShim {
	return &rollupShim{
		logger:         logger,
		clock:          clock,
		windowDuration: windowDuration,
		windows:        map[string]*logWindow{},
	}
}

func (s *rollupShim) WithFields(fields Fields) logShim {
	if len(fields) == 0 {
		return s
	}

	return newRollupShim(
		s.logger.WithFields(fields),
		s.clock,
		s.windowDuration,
	)
}

func (s *rollupShim) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	fields = addCaller(fields)

	if s.getWindow(format).record(s.logger, s.clock, s.windowDuration, level, fields, format, args...) {
		// Not rolling up, log immediately
		// TOOD - explain why layer+2
		s.logger.LogWithFields(level, fields, format, args...)
	}
}

func (s *rollupShim) getWindow(format string) *logWindow {
	s.mutex.RLock()
	if window, ok := s.windows[format]; ok {
		s.mutex.RUnlock()
		return window
	}

	s.mutex.RUnlock()
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if window, ok := s.windows[format]; ok {
		return window
	}

	window := &logWindow{}
	s.windows[format] = window
	return window
}

func (s *rollupShim) Sync() error {
	for _, window := range s.windows {
		window.flush(s.logger)
	}

	return s.logger.Sync()
}

//
// Log Window

func (w *logWindow) record(
	logger Logger,
	clock glock.Clock,
	windowDuration time.Duration,
	level LogLevel,
	fields Fields,
	format string,
	args ...interface{},
) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	now := clock.Now()

	if remaining := windowDuration - now.Sub(w.start); w.start != (time.Time{}) && remaining > 0 {
		w.count++

		if w.count == 1 {
			ch := clock.After(remaining)

			go func() {
				<-ch
				w.flush(logger)
			}()
		}

		return false
	}

	w.flushLocked(logger)

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
	w.mutex.Lock()
	w.flushLocked(logger)
	w.mutex.Unlock()
}

func (w *logWindow) flushLocked(logger Logger) {
	if w.stashed == nil || w.count <= 1 {
		return
	}

	// Set replay field on message
	w.stashed.fields[FieldRollup] = w.count

	logger.LogWithFields(
		w.stashed.level,
		w.stashed.fields,
		w.stashed.format,
		w.stashed.args...,
	)

	w.stashed = nil
}
