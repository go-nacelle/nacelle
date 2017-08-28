package log

import (
	"errors"
	"strings"
)

type Config struct {
	LogBackend       string `env:"LOG_BACKEND" default:"gomol"`
	LogLevel         string `env:"LOG_LEVEL" default:"info"`
	LogEncoding      string `env:"LOG_ENCODING" default:"console"`
	LogDisableCaller bool   `env:"LOG_DISABLE_CALLER"`
	LogInitialFields Fields `env:"LOG_FIELDS"`
}

var (
	ErrIllegalBackend  = errors.New("illegal log backend")
	ErrIllegalLevel    = errors.New("illegal log level")
	ErrIllegalEncoding = errors.New("illegal log encoding")
)

func (c *Config) PostLoad() error {
	c.LogLevel = strings.ToLower(c.LogLevel)

	if !isLegalBackend(c.LogBackend) {
		return ErrIllegalBackend
	}

	if !isLegalLevel(c.LogLevel) {
		return ErrIllegalLevel
	}

	if !isLegalEncoding(c.LogEncoding) {
		return ErrIllegalEncoding
	}

	return nil
}

func isLegalBackend(backend string) bool {
	for _, whitelisted := range []string{"gomol", "logrus", "zap"} {
		if backend == whitelisted {
			return true
		}
	}

	return false
}

func isLegalLevel(level string) bool {
	for _, whitelisted := range []string{"debug", "info", "warning", "error", "fatal"} {
		if level == whitelisted {
			return true
		}
	}

	return false
}

func isLegalEncoding(encoding string) bool {
	return encoding == "console" || encoding == "json"
}
