package nacelle

import (
	configpkg "github.com/go-nacelle/config"
)

type logShim struct {
	logger Logger
}

func (s *logShim) Printf(format string, args ...interface{}) {
	s.logger.Info(format, args...)
}

func NewLoggingConfig(config Config, logger Logger, maskedKeys []string) Config {
	return configpkg.NewLoggingConfig(
		config,
		&logShim{logger},
		maskedKeys,
	)
}
