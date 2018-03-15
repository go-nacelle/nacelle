package nacelle

type (
	// Bootstrapper wraps the entrypoint to the program.
	Bootstrapper struct {
		name            string
		configs         map[interface{}]interface{}
		initFunc        AppInitFunc
		loggingInitFunc LoggingInitFunc
	}

	bootstrapperConfig struct {
		loggingInitFunc LoggingInitFunc
	}

	// AppInitFunc is an program entrypoint called after performing initial
	// configuration loading, sanity checks, and setting up loggers. This
	// function should register initializers and processes and inject values
	// into the service container where necessary.
	AppInitFunc func(*ProcessRunner, *ServiceContainer) error

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

// NewBootstrapper creates an entrypoint to the program with the given configs.
func NewBootstrapper(
	name string,
	configs map[interface{}]interface{},
	initFunc AppInitFunc,
	bootstrapperConfigs ...BoostraperConfigFunc,
) *Bootstrapper {
	if configs == nil {
		configs = map[interface{}]interface{}{}
	}

	config := &bootstrapperConfig{
		loggingInitFunc: InitLogging,
	}

	for _, f := range bootstrapperConfigs {
		f(config)
	}

	return &Bootstrapper{
		name:            name,
		configs:         configs,
		initFunc:        initFunc,
		loggingInitFunc: config.loggingInitFunc,
	}
}

// Boot will initialize services and return a status code - zero
// for a successful exit and one if an error was encountered.
func (bs *Bootstrapper) Boot() int {
	var (
		container = NewServiceContainer()
		runner    = NewProcessRunner(container)
		config    = NewEnvConfig(bs.name)
	)

	for key, obj := range bs.configs {
		if err := config.Register(key, obj); err != nil {
			emergencyLogger().Error("failed to register configs (%s)", err.Error())
			return 1
		}
	}

	if err := config.Register(LoggingConfigToken, &LoggingConfig{}); err != nil {
		emergencyLogger().Error("failed to register logging config (%s)", err.Error())
		return 1
	}

	if errs := config.Load(); len(errs) > 0 {
		logger := emergencyLogger()

		for _, err := range errs {
			logger.Error("Failed to load configuration (%s)", err.Error())
		}

		return 1
	}

	logger, err := bs.loggingInitFunc(config)
	if err != nil {
		emergencyLogger().Error("failed to initialize logging (%s)", err.Error())
		return 1
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			emergencyLogger().Error("failed to sync logs on shutdown (%s)", err.Error())
		}
	}()

	logger.Info("Logging initialized")

	if err := container.Set("logger", logger); err != nil {
		logger.Error("Failed to register logger to service container (%s)", err.Error())
		return 1
	}

	m, err := config.ToMap()
	if err != nil {
		logger.Error("Failed to serialize config (%s)", err.Error())
		return 1
	}

	logger.InfoWithFields(m, "Process starting")

	if err := bs.initFunc(runner, container); err != nil {
		logger.Error("Failed to run initialization function (%s)", err.Error())
		return 1
	}

	statusCode := 0
	for err := range runner.Run(config, logger) {
		statusCode = 1
		logger.Error("Encountered runtime error (%s)", err.Error())
	}

	logger.Info("All processes have stopped")
	return statusCode
}
