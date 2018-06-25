package nacelle

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
