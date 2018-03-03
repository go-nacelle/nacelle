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

func log(logger Logger, level LogLevel, format string, args ...interface{}) {
	switch level {
	case LevelDebug:
		logger.Debug(format, args...)
	case LevelInfo:
		logger.Info(format, args...)
	case LevelWarning:
		logger.Warning(format, args...)
	case LevelError:
		logger.Error(format, args...)
	case LevelFatal:
		logger.Fatal(format, args...)
	}
}

func logWithFields(logger Logger, level LogLevel, fields Fields, format string, args ...interface{}) {
	switch level {
	case LevelDebug:
		logger.DebugWithFields(fields, format, args...)
	case LevelInfo:
		logger.InfoWithFields(fields, format, args...)
	case LevelWarning:
		logger.WarningWithFields(fields, format, args...)
	case LevelError:
		logger.ErrorWithFields(fields, format, args...)
	case LevelFatal:
		logger.FatalWithFields(fields, format, args...)
	}
}
