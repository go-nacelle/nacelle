package log

import (
	"os"

	"github.com/aphistic/gomol"
	console "github.com/aphistic/gomol-console"
)

type GomolShim struct {
	logger gomol.WrappableLogger
}

func NewGomolShim(c *Config) (Logger, error) {
	level, _ := gomol.ToLogLevel(c.LogLevel)
	gomol.SetLogLevel(level)

	if !c.DisableCaller {
		cfg := gomol.NewConfig()
		cfg.FilenameAttr = "filename"
		cfg.LineNumberAttr = "line"
		gomol.SetConfig(cfg)
	}

	if c.LogEncoding == "console" {
		consoleCfg := console.NewConsoleLoggerConfig()
		consoleLogger, _ := console.NewConsoleLogger(consoleCfg)
		consoleLogger.SetTemplate(console.NewTemplateFull())
		gomol.AddLogger(consoleLogger)
	} else {
		// TODO - figure out how to get json into console
	}

	if err := gomol.InitLoggers(); err != nil {
		return nil, err
	}

	return (&GomolShim{logger: gomol.NewLogAdapter(nil)}).WithFields(c.InitialFields), nil
}

func (g *GomolShim) WithFields(fields Fields) Logger {
	if len(fields) == 0 {
		return g
	}

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
	g.log(gomol.LevelFatal, fields, format, args...)
	g.logger.ShutdownLoggers()
	os.Exit(1)
}

func (g *GomolShim) Sync() error {
	return gomol.ShutdownLoggers()
}

func (g *GomolShim) log(level gomol.LogLevel, fields Fields, format string, args ...interface{}) {
	g.logger.Log(level, gomol.NewAttrsFromMap(fields), format, args...)
}
