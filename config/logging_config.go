package config

import (
	"github.com/efritz/nacelle/logging"
	"github.com/efritz/zubrin"
)

type logShim struct {
	logger logging.Logger
}

func (s *logShim) Printf(format string, args ...interface{}) {
	s.logger.Info(format, args...)
}

func NewLoggingConfig(config Config, logger logging.Logger, maskedKeys []string) Config {
	return zubrin.NewLoggingConfig(
		config,
		&logShim{logger},
		maskedKeys,
	)
}
