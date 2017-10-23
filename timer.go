package nacelle

import (
	"time"

	"github.com/efritz/glock"
	"github.com/satori/go.uuid"
)

type (
	Timer struct {
		logger           Logger
		clock            glock.Clock
		id               string
		name             string
		marks            []time.Time
		dropThreshold    time.Duration
		debugThreshold   time.Duration
		infoThreshold    time.Duration
		warningThreshold time.Duration
	}

	TimerConfig func(*Timer)
)

func StartTimer(logger Logger, name string, configs ...TimerConfig) *Timer {
	return startTimer(logger, glock.NewRealClock(), name, configs...)
}

func WithDropThreshold(duration time.Duration) TimerConfig {
	return func(t *Timer) { t.dropThreshold = duration }
}

func WithDebugThreshold(duration time.Duration) TimerConfig {
	return func(t *Timer) { t.debugThreshold = duration }
}

func WithInfoThreshold(duration time.Duration) TimerConfig {
	return func(t *Timer) { t.infoThreshold = duration }
}

func WithWarningThreshold(duration time.Duration) TimerConfig {
	return func(t *Timer) { t.warningThreshold = duration }
}

func startTimer(logger Logger, clock glock.Clock, name string, configs ...TimerConfig) *Timer {
	timer := &Timer{
		logger:           logger,
		clock:            clock,
		id:               uuid.NewV4().String()[:6],
		name:             name,
		marks:            []time.Time{clock.Now()},
		dropThreshold:    time.Millisecond,
		debugThreshold:   time.Millisecond * 200,
		infoThreshold:    time.Millisecond * 600,
		warningThreshold: time.Second,
	}

	for _, config := range configs {
		config(timer)
	}

	attrs := map[string]interface{}{
		"timer_id":   timer.id,
		"timer_name": timer.name,
	}

	logger.InfoWithFields(attrs, "Starting timer: %s", name)
	return timer
}

func (t *Timer) Mark(message string) {
	t.marks = append(t.marks, t.clock.Now())

	var (
		l            = len(t.marks)
		elapsed      = t.marks[l-1].Sub(t.marks[l-2])
		totalElapsed = t.marks[l-1].Sub(t.marks[0])
	)

	logFunc := getLogFunction(
		t.logger,
		elapsed,
		t.dropThreshold,
		t.debugThreshold,
		t.infoThreshold,
		t.warningThreshold,
	)

	attrs := map[string]interface{}{
		"timer_id":           t.id,
		"timer_name":         t.name,
		"elapsed_time":       elapsed,
		"total_elapsed_time": totalElapsed,
	}

	logFunc(attrs, "Marking timer: %s", message)
}

func getLogFunction(logger Logger, elapsed, dropThreshold, debugThreshold, infoThreshold, warningThreshold time.Duration) logFunc {
	if dropThreshold > 0 && elapsed < dropThreshold {
		return noopLogger
	} else if debugThreshold > 0 && elapsed < debugThreshold {
		return logger.DebugWithFields
	} else if infoThreshold > 0 && elapsed < infoThreshold {
		return logger.InfoWithFields
	} else if warningThreshold > 0 && elapsed < warningThreshold {
		return logger.WarningWithFields
	} else {
		return logger.ErrorWithFields
	}
}
