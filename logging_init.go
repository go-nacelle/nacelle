package nacelle

import (
	"fmt"

	"github.com/efritz/nacelle/logging"
)

type loggingConfigToken string

var (
	LoggingConfigToken = loggingConfigToken("nacelle-logging")
	ErrBadConfig       = fmt.Errorf("logging config not registered properly")
)

func InitLogging(config Config) (Logger, error) {
	c := &logging.Config{}
	if err := config.Fetch(LoggingConfigToken, c); err != nil {
		return nil, ErrBadConfig
	}

	return logging.InitGomolShim(c)
}
