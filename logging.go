package nacelle

import (
	"github.com/efritz/nacelle/log"
)

type Logger = log.Logger
type Fields = log.Fields
type LoggingConfig = log.Config

var LoggingConfigToken = struct{}{}

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
