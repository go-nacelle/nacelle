package log

import (
	"errors"
	"strings"
)

type Config struct {
	LogLevel            string `env:"LOG_LEVEL" default:"info"`
	LogEncoding         string `env:"LOG_ENCODING" default:"console"`
	LogColorize         bool   `env:"LOG_COLORIZE" default:"true"`
	LogInitialFields    Fields `env:"LOG_FIELDS"`
	LogShortTime        bool   `env:"LOG_SHORT_TIME" default:"false"`
	LogMultilineFields  bool   `env:"LOG_MULTILINE_FIELDS" default:"true"`
	RawLogAttrBlacklist string `env:"LOG_ATTR_BLACKLIST" default:""`
	LogAttrBlacklist    []string
}

var (
	ErrIllegalLevel    = errors.New("illegal log level")
	ErrIllegalEncoding = errors.New("illegal log encoding")
)

func (c *Config) PostLoad() error {
	c.LogLevel = strings.ToLower(c.LogLevel)

	if !isLegalLevel(c.LogLevel) {
		return ErrIllegalLevel
	}

	if !isLegalEncoding(c.LogEncoding) {
		return ErrIllegalEncoding
	}

	for _, s := range strings.Split(c.RawLogAttrBlacklist, ",") {
		c.LogAttrBlacklist = append(
			c.LogAttrBlacklist,
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
