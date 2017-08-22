package log

import "github.com/sirupsen/logrus"

type LogrusShim struct {
	entry *logrus.Entry
}

func (l *LogrusShim) WithFields(fields Fields) Logger {
	return &LogrusShim{
		entry: l.getEntry(fields),
	}
}

func (l *LogrusShim) Debug(fields Fields, format string, args ...interface{}) {
	l.getEntry(fields).Debugf(format, args...)
}

func (l *LogrusShim) Info(fields Fields, format string, args ...interface{}) {
	l.getEntry(fields).Infof(format, args...)
}

func (l *LogrusShim) Warning(fields Fields, format string, args ...interface{}) {
	l.getEntry(fields).Warningf(format, args...)
}

func (l *LogrusShim) Error(fields Fields, format string, args ...interface{}) {
	l.getEntry(fields).Errorf(format, args...)
}

func (l *LogrusShim) Fatal(fields Fields, format string, args ...interface{}) {
	// TODO - die
	l.getEntry(fields).Fatalf(format, args...)
}

func (l *LogrusShim) getEntry(fields Fields) *logrus.Entry {
	if len(fields) == 0 {
		return l.entry
	}

	return l.entry.WithFields(logrus.Fields(fields))
}
