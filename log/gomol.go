package log

import (
	"os"

	"github.com/aphistic/gomol"
	console "github.com/aphistic/gomol-console"
)

type GomolShim struct {
	logger        gomol.WrappableLogger
	disableCaller bool
}

func NewGomolShim(c *Config) (Logger, error) {
	level, _ := gomol.ToLogLevel(c.LogLevel)
	gomol.SetLogLevel(level)

	if c.LogEncoding == "console" {
		consoleLogger, err := console.NewConsoleLogger(&console.ConsoleLoggerConfig{
			Colorize: true,
			Writer:   os.Stderr,
		})

		if err != nil {
			return nil, err
		}

		tpl, err := gomol.NewTemplate("" +
			"[{{.Timestamp.Format \"2006-01-02 15:04:05.000\"}}] " +
			`{{color}}{{printf "%5s" (ucase .LevelName)}}{{reset}} ` +
			"{{.Message}}" +
			"{{if .Attrs}}{{range $key, $val := .Attrs}} {{color}}{{$key}}{{reset}}={{$val}}{{end}}{{end}}")

		if err != nil {
			return nil, err
		}

		consoleLogger.SetTemplate(tpl)
		gomol.AddLogger(consoleLogger)
	} else {
		gomol.AddLogger(newJSONLogger())
	}

	if err := gomol.InitLoggers(); err != nil {
		return nil, err
	}

	shim := &GomolShim{
		logger:        gomol.NewLogAdapter(nil),
		disableCaller: c.LogDisableCaller,
	}

	return shim.WithFields(c.LogInitialFields), nil
}

func (g *GomolShim) WithFields(fields Fields) Logger {
	if len(fields) == 0 {
		return g
	}

	return &GomolShim{
		logger:        gomol.NewLogAdapterFor(g.logger, gomol.NewAttrsFromMap(fields)),
		disableCaller: g.disableCaller,
	}
}

func (g *GomolShim) Debug(format string, args ...interface{}) {
	g.log(gomol.LevelDebug, nil, format, args...)
}

func (g *GomolShim) Info(format string, args ...interface{}) {
	g.log(gomol.LevelInfo, nil, format, args...)
}

func (g *GomolShim) Warning(format string, args ...interface{}) {
	g.log(gomol.LevelWarning, nil, format, args...)
}

func (g *GomolShim) Error(format string, args ...interface{}) {
	g.log(gomol.LevelError, nil, format, args...)
}

func (g *GomolShim) Fatal(format string, args ...interface{}) {
	g.log(gomol.LevelFatal, nil, format, args...)
	g.logger.ShutdownLoggers()
	os.Exit(1)
}

func (g *GomolShim) DebugWithFields(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelDebug, fields, format, args...)
}

func (g *GomolShim) InfoWithFields(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelInfo, fields, format, args...)
}

func (g *GomolShim) WarningWithFields(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelWarning, fields, format, args...)
}

func (g *GomolShim) ErrorWithFields(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelError, fields, format, args...)
}

func (g *GomolShim) FatalWithFields(fields Fields, format string, args ...interface{}) {
	g.log(gomol.LevelFatal, fields, format, args...)
	g.logger.ShutdownLoggers()
	os.Exit(1)
}

func (g *GomolShim) Sync() error {
	return gomol.ShutdownLoggers()
}

func (g *GomolShim) log(level gomol.LogLevel, fields Fields, format string, args ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}

	if !g.disableCaller {
		fields["caller"] = getCaller()
	}

	g.logger.Log(level, gomol.NewAttrsFromMap(fields.normalizeTimeValues()), format, args...)
}
