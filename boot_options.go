package nacelle

import (
	"fmt"

	"github.com/efritz/nacelle/logging"
)

type (
	// LoggingInitFunc creates a factory from a config object.
	LoggingInitFunc func(Config) (Logger, error)

	// BootstrapperConfigFunc is a function used to configure an instance of
	// a Bootstrapper.
	BootstrapperConfigFunc func(*bootstrapperConfig)
)

// WithLoggingInitFunc sets the function that initializes logging.
func WithLoggingInitFunc(loggingInitFunc LoggingInitFunc) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.loggingInitFunc = loggingInitFunc }
}

// WithLoggingFields sets additional fields sent with every log message.
func WithLoggingFields(loggingFields Fields) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.loggingFields = loggingFields }
}

// WithRunnerOptions passes RunnerConfigFuncs to the runner created by Boot.
func WithRunnerOptions(configs ...RunnerConfigFunc) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.runnerConfigFuncs = configs }
}

func defaultLogginInitFunc(config Config) (Logger, error) {
	c := &logging.Config{}
	if err := config.Fetch(loggingConfigToken, c); err != nil {
		return nil, fmt.Errorf("logging config not registered properly")
	}

	return logging.InitGomolShim(c)
}
