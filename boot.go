package nacelle

import (
	"os"

	"github.com/efritz/nacelle/config"
	"github.com/efritz/nacelle/logging"
	"github.com/efritz/nacelle/process"
	"github.com/efritz/nacelle/service"
)

type (
	// Bootstrapper wraps the entrypoint to the program.
	Bootstrapper struct {
		configs           map[interface{}]interface{}
		initFunc          AppInitFunc
		configSourcer     ConfigSourcer
		loggingInitFunc   LoggingInitFunc
		loggingFields     LogFields
		runnerConfigFuncs []RunnerConfigFunc
	}

	bootstrapperConfig struct {
		configSourcer     ConfigSourcer
		loggingInitFunc   LoggingInitFunc
		loggingFields     LogFields
		runnerConfigFuncs []RunnerConfigFunc
	}

	// AppInitFunc is an program entrypoint called after performing initial
	// configuration loading, sanity checks, and setting up loggers. This
	// function should register initializers and processes and inject values
	// into the service container where necessary.
	AppInitFunc func(ProcessContainer, ServiceContainer) error
)

// NewBootstrapper creates an entrypoint to the program with the given configs.
func NewBootstrapper(
	name string,
	initFunc AppInitFunc,
	bootstrapperConfigs ...BootstrapperConfigFunc,
) *Bootstrapper {
	config := &bootstrapperConfig{
		configSourcer:   NewEnvSourcer(name),
		loggingInitFunc: defaultLogginInitFunc,
	}

	for _, f := range bootstrapperConfigs {
		f(config)
	}

	return &Bootstrapper{
		initFunc:          initFunc,
		configSourcer:     config.configSourcer,
		loggingInitFunc:   config.loggingInitFunc,
		loggingFields:     config.loggingFields,
		runnerConfigFuncs: config.runnerConfigFuncs,
	}
}

// Boot will initialize services and return a status code. This
// method does not return in any meaningful way (it blocks until
// the associated process runner has completed).
func (bs *Bootstrapper) Boot() int {
	baseConfig := NewConfig(bs.configSourcer)

	logger, err := bs.loggingInitFunc(baseConfig)
	if err != nil {
		logging.LogEmergencyError("failed to initialize logging (%s)", err)
		return 1
	}

	logger = logger.WithFields(bs.loggingFields)

	defer func() {
		if err := logger.Sync(); err != nil {
			logging.LogEmergencyError("failed to sync logs on shutdown (%s)", err)
		}
	}()

	logger.Info("Logging initialized")

	serviceContainer, err := service.NewContainer()
	if err != nil {
		logging.LogEmergencyError("failed to create service container (%s)", err)
		return 1
	}

	if err := serviceContainer.Set("logger", logger); err != nil {
		logger.Error("Failed to register logger to service container (%s)", err)
		return 1
	}

	health := process.NewHealth()

	if err := serviceContainer.Set("health", health); err != nil {
		logger.Error("Failed to register health reporter to service container (%s)", err)
		return 1
	}

	config := config.NewLoggingConfig(baseConfig, logger)
	processContainer := process.NewContainer()

	if err := bs.initFunc(processContainer, serviceContainer); err != nil {
		logger.Error("Failed to run initialization function (%s)", err.Error())
		return 1
	}

	runner := process.NewRunner(
		processContainer,
		serviceContainer,
		health,
		append(
			bs.runnerConfigFuncs,
			process.WithLogger(logger),
		)...,
	)

	statusCode := 0
	for range runner.Run(config) {
		statusCode = 1
	}

	logger.Info("All processes have stopped")
	return statusCode
}

// BootAndExit calls Boot and sets the program return code on halt. This
// method does not return.
func (bs *Bootstrapper) BootAndExit() {
	os.Exit(bs.Boot())
}
