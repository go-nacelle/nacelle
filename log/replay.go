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
	// ReplayLogger is a Logger that provides a way to replay a sequence of
	// message in the order they were logged, at a higher log level.
	ReplayLogger interface {
		Logger

		// Replay will cause all of the messages previously logged at one of the
		// journaled levels to be re-set at the given level. All future messages
		// logged at one of the journaled levels will be replayed immediately.
		Replay(LogLevel)
	}

	replayShim struct {
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

var _ logShim = &replayShim{}

// NewReplayAdapter creates a ReplayLogger wrapping the given logger.
func NewReplayAdapter(logger Logger, levels ...LogLevel) ReplayLogger {
	return adaptReplayShim(newReplayShim(logger, glock.NewRealClock(), levels...))
}

func newReplayShim(logger Logger, clock glock.Clock, levels ...LogLevel) *replayShim {
	sharedJournal := &sharedJournal{
		clock:    clock,
		messages: []*journaledMessage{},
		levels:   levels,
	}

	return &replayShim{
		logger:        logger,
		sharedJournal: sharedJournal,
	}
}

func (s *replayShim) WithFields(fields Fields) logShim {
	if len(fields) == 0 {
		return s
	}

	return &replayShim{
		logger:        s.logger.WithFields(fields),
		sharedJournal: s.sharedJournal,
	}
}

func (s *replayShim) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	fields = addCaller(fields)

	// Log immediately
	s.logger.LogWithFields(level, fields, format, args...)

	// Add to journal
	s.sharedJournal.record(s.logger, level, fields, format, args)
}

func (s *replayShim) Sync() error {
	return s.logger.Sync()
}

func (s *replayShim) Replay(level LogLevel) {
	s.sharedJournal.replay(level)
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

	m.logger.LogWithFields(
		*level,
		m.message.fields,
		m.message.format,
		m.message.args...,
	)
}
