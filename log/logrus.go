package log

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type LogrusShim struct {
	entry *logrus.Entry
}

//
// Shim

func NewLogrusLogger(logger *logrus.Entry, initialFields Fields) Logger {
	return adaptShim((&LogrusShim{logger}).WithFields(initialFields))
}

func (l *LogrusShim) WithFields(fields Fields) logShim {
	if len(fields) == 0 {
		return l
	}

	return &LogrusShim{l.getEntry(fields)}
}

func (l *LogrusShim) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	entry := l.getEntry(addCaller(fields))

	switch level {
	case LevelDebug:
		entry.Debugf(format, args...)
	case LevelInfo:
		entry.Infof(format, args...)
	case LevelWarning:
		entry.Warningf(format, args...)
	case LevelError:
		entry.Errorf(format, args...)
	case LevelFatal:
		entry.Fatalf(format, args...)
	}
}

func (l *LogrusShim) Sync() error {
	return nil
}

func (l *LogrusShim) getEntry(fields Fields) *logrus.Entry {
	if len(fields) == 0 {
		return l.entry
	}

	return l.entry.WithFields(logrus.Fields(fields.normalizeTimeValues()))
}

//
// Init

func InitLogrusShim(c *Config) (Logger, error) {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.Level = level

	if c.LogEncoding == "console" {
		formatter := &prefixed.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  ConsoleTimeFormat,
			QuoteEmptyFields: true,
			ForceColors:      c.LogColorize,
		}

		logger.Formatter = formatter
	} else {
		logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: JSONTimeFormat,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		}
	}

	return NewLogrusLogger(logger.WithFields(nil), c.LogInitialFields), nil
}
