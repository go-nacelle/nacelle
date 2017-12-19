package nacelle

import (
	"errors"

	"github.com/efritz/nacelle/log"
)

type (
	Logger        = log.Logger
	Fields        = log.Fields
	LoggingConfig = log.Config

	loggingConfigToken struct{}
	logFunc            func(log.Fields, string, ...interface{})
)

var (
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
		LogLevel:         "DEBUG",
		LogEncoding:      "json",
		LogDisableCaller: true,
	})

	return logger
}

func noopLogger(fields log.Fields, message string, args ...interface{}) {
	// Silence is golden
}
