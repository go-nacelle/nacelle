package log

type (
	Logger interface {
		WithFields(fields Fields) Logger
		Debug(fields Fields, format string, args ...interface{})
		Info(fields Fields, format string, args ...interface{})
		Warning(fields Fields, format string, args ...interface{})
		Error(fields Fields, format string, args ...interface{})
		Fatal(fields Fields, format string, args ...interface{})
	}

	Fields map[string]interface{}
)
