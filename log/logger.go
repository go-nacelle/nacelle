package log

type (
	Logger interface {
		WithFields(Fields) Logger
		LogWithFields(LogLevel, Fields, string, ...interface{})
		Sync() error

		// Convenience Methods
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
	}

	LogLevel int
)

const (
	LevelFatal LogLevel = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug

	ConsoleTimeFormat = "2006-01-02 15:04:05.000"
	JSONTimeFormat    = "2006-01-02T15:04:05.000-0700"
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}
