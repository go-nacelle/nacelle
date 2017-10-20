package log

import "time"

type (
	Logger interface {
		WithFields(fields Fields) Logger
		Debug(format string, args ...interface{})
		Info(format string, args ...interface{})
		Warning(format string, args ...interface{})
		Error(format string, args ...interface{})
		Fatal(format string, args ...interface{})
		DebugWithFields(fields Fields, format string, args ...interface{})
		InfoWithFields(fields Fields, format string, args ...interface{})
		WarningWithFields(fields Fields, format string, args ...interface{})
		ErrorWithFields(fields Fields, format string, args ...interface{})
		FatalWithFields(fields Fields, format string, args ...interface{})
		Sync() error
	}

	Fields map[string]interface{}
)

const (
	ConsoleTimeFormat = "2006-01-02 15:04:05.000"
	JSONTimeFormat    = "2006-01-02T15:04:05.000-0700"
)

func (f Fields) normalizeTimeValues() Fields {
	for key, val := range f {
		switch v := val.(type) {
		case time.Time:
			f[key] = v.Format(JSONTimeFormat)
		default:
			f[key] = v
		}
	}

	return f
}
