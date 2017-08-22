package log

import (
	"errors"
	"strings"
)

type (
	Logger interface {
		WithFields(fields Fields) Logger
		Debug(fields Fields, format string, args ...interface{})
		Info(fields Fields, format string, args ...interface{})
		Warning(fields Fields, format string, args ...interface{})
		Error(fields Fields, format string, args ...interface{})
		Fatal(fields Fields, format string, args ...interface{})
		Sync() error
	}

	Fields map[string]interface{}

	Config struct {
		LogLevel      string `env:"LOG_LEVEL" default:"info"`
		LogEncoding   string `env:"LOG_ENCODING" default:"console"`
		DisableCaller bool   `env:"LOG_ENABLE_CALLER"`
		InitialFields Fields `env:"LOG_FIELDS"`
	}
)

var (
	ErrIllegalLevel    = errors.New("illegal log level")
	ErrIllegalEncoding = errors.New("illegal log encoding")
)

func (c *Config) PostLoad() error {
	c.LogLevel = strings.ToLower(c.LogLevel)

	if !isLegalLevel(c.LogLevel) {
		return ErrIllegalLevel
	}

	if c.LogEncoding != "console" && c.LogEncoding != "json" {
		return ErrIllegalEncoding
	}

	return nil
}

func isLegalLevel(level string) bool {
	for _, whitelisted := range []string{"debug", "info", "warning", "error", "fatal"} {
		if level == whitelisted {
			return true
		}
	}

	return false
}
