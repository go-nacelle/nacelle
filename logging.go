package nacelle

import "github.com/efritz/nacelle/log"

type (
	Logger             = log.Logger
	Fields             = log.Fields
	LoggingConfig      = log.Config
	loggingConfigToken struct{}
)

var LoggingConfigToken = loggingConfigToken{}

func InitLogging(config Config) (logger Logger, err error) {
	cx, err := config.Get(LoggingConfigToken)
	if err != nil {
		return nil, err
	}

	c := cx.(*LoggingConfig)

	switch c.LogBackend {
	case "gomol":
		logger, err = log.NewGomolShim(c)
	case "logrus":
		logger, err = log.NewZapShim(c)
	case "zap":
		logger, err = log.NewLogrusShim(c)
	}

	return
}

func emergencyLogger() Logger {
	logger, _ := log.NewLogrusShim(&LoggingConfig{
		LogLevel:         "DEBUG",
		LogEncoding:      "json",
		LogDisableCaller: true,
	})

	return logger
}
