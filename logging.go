package nacelle

import (
	"errors"

	"github.com/efritz/nacelle/log"
)

type (
	Logger        = log.Logger
	ReplayLogger  = log.ReplayLogger
	Fields        = log.Fields
	LoggingConfig = log.Config
	LogLevel      = log.LogLevel

	loggingConfigToken struct{}
	logFunc            func(log.Fields, string, ...interface{})
)

const (
	LevelFatal   = log.LevelFatal
	LevelError   = log.LevelError
	LevelWarning = log.LevelWarning
	LevelInfo    = log.LevelInfo
	LevelDebug   = log.LevelDebug
)

var (
	NewReplayAdapter = log.NewReplayAdapter
	NewRollupAdapter = log.NewRollupAdapter

	LoggingConfigToken = loggingConfigToken{}
	ErrBadConfig       = errors.New("logging config not registered properly")
)

func InitLogging(config Config) (logger Logger, err error) {
	cx, err := config.Get(LoggingConfigToken)
	if err != nil {
		return nil, err
	}

	c, ok := cx.(*LoggingConfig)
	if !ok {
		return nil, ErrBadConfig
	}

	switch c.LogBackend {
	case "gomol":
		logger, err = log.InitGomolShim(c)
	case "logrus":
		logger, err = log.InitZapShim(c)
	case "zap":
		logger, err = log.InitLogrusShim(c)
	}

	return
}

func emergencyLogger() Logger {
	logger, _ := log.InitLogrusShim(&LoggingConfig{
		LogLevel:    "DEBUG",
		LogEncoding: "json",
	})

	return logger
}

func noopLogger(fields log.Fields, message string, args ...interface{}) {
	// Silence is golden
}
