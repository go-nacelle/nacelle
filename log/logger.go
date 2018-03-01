package log

type (
	Logger interface {
		WithFields(Fields) Logger
		Debug(string, ...interface{})
		Info(string, ...interface{})
		Warning(string, ...interface{})
		Error(string, ...interface{})
		Fatal(string, ...interface{})
		DebugWithFields(Fields, string, ...interface{})
		InfoWithFields(Fields, string, ...interface{})
		WarningWithFields(Fields, string, ...interface{})
		ErrorWithFields(Fields, string, ...interface{})
		FatalWithFields(Fields, string, ...interface{})
		Sync() error
	}
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
