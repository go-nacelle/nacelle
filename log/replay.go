package log

import (
	"sync"

	"github.com/efritz/glock"
)

// FieldReplay is a field assigned to a message that has
// been replayed at a different log level. Its value is equal
// to the original log level.
const FieldReplay = "replayed-from-level"

type (
	// ReplayAdapter provides a way to replay a sequence of message, in
	// the order they were logged, at a higher log level.
	ReplayAdapter struct {
		logger        Logger
		sharedJournal *sharedJournal
	}

	sharedJournal struct {
		clock       glock.Clock
		messages    []*journaledMessage
		levels      []LogLevel
		replayingAt *LogLevel
		mutex       sync.RWMutex
	}

	journaledMessage struct {
		logger  Logger
		message logMessage
	}
)

//
// Shim

var _ logShim = &ReplayAdapter{}

// NewReplayAdapter creates a ReplayAdapter which wraps the given logger.
func NewReplayAdapter(logger Logger, levels ...LogLevel) *ReplayAdapter {
	return newReplayAdapter(logger, glock.NewRealClock(), levels...)
}

func newReplayAdapter(logger Logger, clock glock.Clock, levels ...LogLevel) *ReplayAdapter {
	sharedJournal := &sharedJournal{
		clock:    clock,
		messages: []*journaledMessage{},
		levels:   levels,
	}

	return &ReplayAdapter{
		logger:        logger,
		sharedJournal: sharedJournal,
	}
}

func (a *ReplayAdapter) WithFields(fields Fields) logShim {
	if len(fields) == 0 {
		return a
	}

	return &ReplayAdapter{
		logger:        a.logger.WithFields(fields),
		sharedJournal: a.sharedJournal,
	}
}

func (a *ReplayAdapter) Log(level LogLevel, format string, args ...interface{}) {
	a.LogWithFields(level, nil, format, args)
}

func (a *ReplayAdapter) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	// Log immediately
	logWithFields(a.logger, level, fields, format, args...)

	// Add to journal
	a.sharedJournal.record(a.logger, level, fields, format, args)
}

func (a *ReplayAdapter) Sync() error {
	return a.logger.Sync()
}

// TODO - how to export properly? This is a shim.

// Replay will cause all of the messages previously logged at one of the
// journaled levels to be re-set at the given level. All future messages
// logged at one of the journaled levels will be replayed immediately.
func (a *ReplayAdapter) Replay(level LogLevel) {
	a.sharedJournal.replay(level)
}

//
// Shared Journal

func (j *sharedJournal) record(logger Logger, level LogLevel, fields Fields, format string, args []interface{}) {
	if !j.shouldJournal(level) {
		return
	}

	innerMessage := logMessage{
		level:  level,
		fields: fields.clone(),
		format: format,
		args:   args,
	}

	message := &journaledMessage{
		logger:  logger,
		message: innerMessage,
	}

	j.mutex.RLock()
	message.replay(j.replayingAt)
	j.mutex.RUnlock()

	j.mutex.Lock()
	j.messages = append(j.messages, message)
	j.mutex.Unlock()
}

func (j *sharedJournal) shouldJournal(level LogLevel) bool {
	for _, l := range j.levels {
		if l == level {
			return true
		}
	}

	return false
}

func (j *sharedJournal) replay(level LogLevel) {
	j.mutex.RLock()
	shouldReplay := j.replayingAt == nil || level < *j.replayingAt
	j.mutex.RUnlock()

	if !shouldReplay {
		return
	}

	j.mutex.Lock()
	j.replayingAt = &level
	j.mutex.Unlock()

	j.mutex.RLock()
	defer j.mutex.RUnlock()

	for _, message := range j.messages {
		message.replay(&level)
	}
}

func (m *journaledMessage) replay(level *LogLevel) {
	if level == nil {
		return
	}

	// Set replay field on message
	m.message.fields[FieldReplay] = m.message.level

	logWithFields(
		m.logger,
		*level,
		m.message.fields,
		m.message.format,
		m.message.args...,
	)
}
