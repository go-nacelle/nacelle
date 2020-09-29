package nacelle

import "github.com/go-nacelle/config"

type logShim struct {
	logger Logger
}

var newLoggingConfig = config.NewLoggingConfig

func (s *logShim) Printf(format string, args ...interface{}) {
	s.logger.Info(format, args...)
}

func NewLoggingConfig(config Config, logger Logger, maskedKeys []string) Config {
	return newLoggingConfig(
		config,
		&logShim{logger},
		maskedKeys,
	)
}
