package logging

import (
	"fmt"
	"strings"
)

type Config struct {
	LogLevel                  string `env:"LOG_LEVEL" default:"info"`
	LogEncoding               string `env:"LOG_ENCODING" default:"console"`
	LogColorize               bool   `env:"LOG_COLORIZE" default:"true"`
	LogInitialFields          Fields `env:"LOG_FIELDS"`
	LogShortTime              bool   `env:"LOG_SHORT_TIME" default:"false"`
	LogDisplayFields          bool   `env:"LOG_DISPLAY_FIELDS" default:"true"`
	LogDisplayMultilineFields bool   `env:"LOG_DISPLAY_MULTILINE_FIELDS" default:"false"`
	RawLogFieldBlacklist      string `env:"LOG_FIELD_BLACKLIST"`
	LogFieldBlacklist         []string
}

var (
	ErrIllegalLevel    = fmt.Errorf("illegal log level")
	ErrIllegalEncoding = fmt.Errorf("illegal log encoding")
)

func (c *Config) PostLoad() error {
	c.LogLevel = strings.ToLower(c.LogLevel)

	if !isLegalLevel(c.LogLevel) {
		return ErrIllegalLevel
	}

	if !isLegalEncoding(c.LogEncoding) {
		return ErrIllegalEncoding
	}

	for _, s := range strings.Split(c.RawLogFieldBlacklist, ",") {
		c.LogFieldBlacklist = append(
			c.LogFieldBlacklist,
			strings.ToLower(strings.TrimSpace(s)),
		)
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

func isLegalEncoding(encoding string) bool {
	return encoding == "console" || encoding == "json"
}
