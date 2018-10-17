package logging

import (
	"fmt"
	"strings"
)

type Config struct {
	LogLevel                  string `env:"log_level" file:"log_level" default:"info"`
	LogEncoding               string `env:"log_encoding" file:"log_encoding" default:"console"`
	LogColorize               bool   `env:"log_colorize" file:"log_colorize" default:"true"`
	LogInitialFields          Fields `env:"log_fields" file:"log_fields"`
	LogShortTime              bool   `env:"log_short_time" file:"log_short_time" default:"false"`
	LogDisplayFields          bool   `env:"log_display_fields" file:"log_display_fields" default:"true"`
	LogDisplayMultilineFields bool   `env:"log_display_multiline_fields" file:"log_display_multiline_fields" default:"false"`
	RawLogFieldBlacklist      string `env:"log_field_blacklist" file:"log_field_blacklist" mask:"true"`
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
