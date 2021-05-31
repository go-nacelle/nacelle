package nacelle

import (
	"context"
	"fmt"

	"github.com/go-nacelle/log/v2"
)

// LoggingInitFunc creates a factory from a config object.
type LoggingInitFunc func(*Config) (Logger, error)

// BootstrapperConfigFunc is a function used to configure an instance of
// a Bootstrapper.
type BootstrapperConfigFunc func(*bootstrapperConfig)

// WithContextFilter sets the context filter for the bootstrapper.
func WithContextFilter(f func(ctx context.Context) context.Context) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.contextFilter = f }
}

// WithConfigSourcer sets the source that should be used for populating config structs.
func WithConfigSourcer(configSourcer ConfigSourcer) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.configSourcer = configSourcer }
}

// WithConfigMaskedKeys sets the keys that are redacted when printed by the config logger.
func WithConfigMaskedKeys(configMaskedKeys []string) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.configMaskedKeys = configMaskedKeys }
}

// WithLoggingInitFunc sets the function that initializes logging.
func WithLoggingInitFunc(loggingInitFunc LoggingInitFunc) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.loggingInitFunc = loggingInitFunc }
}

// WithLoggingFields sets additional fields sent with every log message.
func WithLoggingFields(loggingFields LogFields) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.loggingFields = loggingFields }
}

// WithMachineOptions applies additional configuration to the process machine.
func WithMachineOptions(configs ...MachineConfigFunc) BootstrapperConfigFunc {
	return func(c *bootstrapperConfig) { c.machineConfigFuncs = configs }
}

func defaultLoggingInitFunc(config *Config) (Logger, error) {
	c := &log.Config{}
	if err := config.Load(c); err != nil {
		return nil, fmt.Errorf("could not load logging config (%s)", err.Error())
	}

	return log.InitLogger(c)
}
