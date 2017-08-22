package log

import (
	"github.com/aphistic/gomol"
	console "github.com/aphistic/gomol-console"
)

type GomolShim struct {
	logger gomol.WrappableLogger
}

func (g *GomolShim) WithFields(fields Fields) Logger {
	return &GomolShim{
		logger: gomol.NewLogAdapterFor(g.logger, gomol.NewAttrsFromMap(fields)),
	}
}

func (g *GomolShim) Debug(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelDebug, fields, format, args...)
}

func (g *GomolShim) Info(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelInfo, fields, format, args...)
}

func (g *GomolShim) Warning(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelWarning, fields, format, args...)
}

func (g *GomolShim) Error(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelError, fields, format, args...)
}

func (g *GomolShim) Fatal(fields Fields, format string, args ...interface{}) {
	// TODO - die
	g.log(gomol.LevelFatal, fields, format, args...)
}

func (g *GomolShim) log(level gomol.LogLevel, fields Fields, format string, args ...interface{}) {
	g.logger.Log(level, gomol.NewAttrsFromMap(fields), format, args...)
}
