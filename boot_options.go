package nacelle

type (
	// LoggingInitFunc creates a factory from a config object.
	LoggingInitFunc func(Config) (Logger, error)

	// BoostraperConfigFunc is a function used to configure an instance of
	// a Bootstrapper.
	BoostraperConfigFunc func(*bootstrapperConfig)
)

// WithLoggingInitFunc sets the function that initializes logging.
func WithLoggingInitFunc(loggingInitFunc LoggingInitFunc) BoostraperConfigFunc {
	return func(c *bootstrapperConfig) { c.loggingInitFunc = loggingInitFunc }
}
