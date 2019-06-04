package config

import (
	configlogging "github.com/go-nacelle/config/logging"
	"github.com/go-nacelle/nacelle/logging"
)

type logShim struct {
	logger logging.Logger
}

func (s *logShim) Printf(format string, args ...interface{}) {
	s.logger.Info(format, args...)
}

func NewLoggingConfig(config Config, logger logging.Logger, maskedKeys []string) Config {
	return configlogging.NewLoggingConfig(
		config,
		&logShim{logger},
		maskedKeys,
	)
}
