package log

type (
	logShim interface {
		WithFields(Fields) logShim
		Log(LogLevel, string, ...interface{})
		LogWithFields(LogLevel, Fields, string, ...interface{})
		Sync() error
	}

	shimAdapter struct {
		shim logShim
	}

	logMessage struct {
		level  LogLevel
		fields Fields
		format string
		args   []interface{}
	}
)

func adaptShim(shim logShim) Logger {
	return &shimAdapter{shim: shim}
}

func (sa *shimAdapter) WithFields(fields Fields) Logger {
	if len(fields) == 0 {
		return sa
	}

	return &shimAdapter{shim: sa.shim.WithFields(fields)}
}

func (sa *shimAdapter) Debug(format string, args ...interface{}) {
	sa.shim.Log(LevelDebug, format, args)
}

func (sa *shimAdapter) Info(format string, args ...interface{}) {
	sa.shim.Log(LevelInfo, format, args)
}

func (sa *shimAdapter) Warning(format string, args ...interface{}) {
	sa.shim.Log(LevelWarning, format, args)
}

func (sa *shimAdapter) Error(format string, args ...interface{}) {
	sa.shim.Log(LevelError, format, args)
}

func (sa *shimAdapter) Fatal(format string, args ...interface{}) {
	sa.shim.Log(LevelFatal, format, args)
}

func (sa *shimAdapter) DebugWithFields(fields Fields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelDebug, fields, format, args)
}

func (sa *shimAdapter) InfoWithFields(fields Fields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelInfo, fields, format, args)
}

func (sa *shimAdapter) WarningWithFields(fields Fields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelWarning, fields, format, args)
}

func (sa *shimAdapter) ErrorWithFields(fields Fields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelError, fields, format, args)
}

func (sa *shimAdapter) FatalWithFields(fields Fields, format string, args ...interface{}) {
	sa.shim.LogWithFields(LevelFatal, fields, format, args)
}

func (sa *shimAdapter) Sync() error {
	return sa.shim.Sync()
}
