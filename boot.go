package nacelle

import (
	"os"
	// TODO - get rid of all imports soutside of imports.go
	"github.com/go-nacelle/process"
)

type (
	// Bootstrapper wraps the entrypoint to the program.
	Bootstrapper struct {
		configs           map[interface{}]interface{}
		initFunc          AppInitFunc
		configSourcer     ConfigSourcer
		configMaskedKeys  []string
		loggingInitFunc   LoggingInitFunc
		loggingFields     LogFields
		runnerConfigFuncs []RunnerConfigFunc
	}

	bootstrapperConfig struct {
		configSourcer     ConfigSourcer
		configMaskedKeys  []string
		loggingInitFunc   LoggingInitFunc
		loggingFields     LogFields
		runnerConfigFuncs []RunnerConfigFunc
	}

	// AppInitFunc is an program entrypoint called after performing initial
	// configuration loading, sanity checks, and setting up loggers. This
	// function should register initializers and processes and inject values
	// into the service container where necessary.
	AppInitFunc func(ProcessContainer, ServiceContainer) error

	// ServiceInitializerFunc is an InitializerFunc with a service container argument.
	ServiceInitializerFunc func(config Config, container ServiceContainer) error
)

// WrapServiceInitializerFunc creates an InitializerFunc from a ServiceInitializerFunc and a container.
func WrapServiceInitializerFunc(container ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
	return InitializerFunc(func(config Config) error {
		return f(config, container)
	})
}

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
		configMaskedKeys:  config.configMaskedKeys,
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
		LogEmergencyError("failed to initialize logging (%s)", err)
		return 1
	}

	logger = logger.WithFields(bs.loggingFields)

	defer func() {
		if err := logger.Sync(); err != nil {
			LogEmergencyError("failed to sync logs on shutdown (%s)", err)
		}
	}()

	logger.Info("Logging initialized")

	health := NewHealth()

	serviceContainer := NewServiceContainer()
	_ = serviceContainer.Set("health", health)
	_ = serviceContainer.Set("logger", logger)
	_ = serviceContainer.Set("services", serviceContainer)

	config := NewLoggingConfig(baseConfig, logger, bs.configMaskedKeys)
	processContainer := NewProcessContainer()

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
