package log

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type LogrusShim struct {
	entry         *logrus.Entry
	disableCaller bool
}

func NewLogrusShim(c *Config) (Logger, error) {
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

	shim := &LogrusShim{
		entry:         logger.WithFields(nil),
		disableCaller: c.LogDisableCaller,
	}

	return shim.WithFields(c.LogInitialFields), nil
}

func (l *LogrusShim) WithFields(fields Fields) Logger {
	if len(fields) == 0 {
		return l
	}

	return &LogrusShim{
		entry:         l.getEntry(fields),
		disableCaller: l.disableCaller,
	}
}

func (l *LogrusShim) Debug(fields Fields, format string, args ...interface{}) {
	l.decorateEntry(fields).Debugf(format, args...)
}

func (l *LogrusShim) Info(fields Fields, format string, args ...interface{}) {
	l.decorateEntry(fields).Infof(format, args...)
}

func (l *LogrusShim) Warning(fields Fields, format string, args ...interface{}) {
	l.decorateEntry(fields).Warningf(format, args...)
}

func (l *LogrusShim) Error(fields Fields, format string, args ...interface{}) {
	l.decorateEntry(fields).Errorf(format, args...)
}

func (l *LogrusShim) Fatal(fields Fields, format string, args ...interface{}) {
	l.decorateEntry(fields).Fatalf(format, args...)
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

func (l *LogrusShim) decorateEntry(fields Fields) *logrus.Entry {
	entry := l.getEntry(fields)
	if l.disableCaller {
		return entry
	}

	return entry.WithFields(logrus.Fields(map[string]interface{}{
		"caller": getCaller(),
	}))
}
