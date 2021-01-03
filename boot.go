package nacelle

import (
	"context"
	"fmt"
	"os"

	"github.com/go-nacelle/config"
	"github.com/go-nacelle/log"
	"github.com/go-nacelle/process"
)

// Bootstrapper wraps the entrypoint to the program.
type Bootstrapper struct {
	initFunc          AppInitFunc
	contextFilter     func(ctx context.Context) context.Context
	configSourcer     ConfigSourcer
	configMaskedKeys  []string
	loggingInitFunc   LoggingInitFunc
	loggingFields     LogFields
	runnerConfigFuncs []RunnerConfigFunc
}

type bootstrapperConfig struct {
	contextFilter     func(ctx context.Context) context.Context
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
type AppInitFunc func(ProcessContainer, ServiceContainer) error

// ServiceInitializerFunc is an InitializerFunc with a service container argument.
type ServiceInitializerFunc func(ctx context.Context, container ServiceContainer) error

// WrapServiceInitializerFunc creates an InitializerFunc from a ServiceInitializerFunc and a container.
func WrapServiceInitializerFunc(container ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
	return InitializerFunc(func(ctx context.Context) error {
		return f(ctx, container)
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
		loggingInitFunc: defaultLoggingInitFunc,
	}

	for _, f := range bootstrapperConfigs {
		f(config)
	}

	return &Bootstrapper{
		initFunc:          initFunc,
		contextFilter:     config.contextFilter,
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
	showHelp := showHelp()

	shim := &logShim{}
	config := NewConfig(
		bs.configSourcer,
		config.WithLogger(shim),
		config.WithMaskedKeys(bs.configMaskedKeys),
	)

	if err := config.Init(); err != nil {
		LogEmergencyError("failed to initialize config (%s)", err)
		return 1
	}

	logger, err := bs.makeLogger(config, !showHelp)
	if err != nil {
		LogEmergencyError("failed to initialize logging (%s)", err)
		return 1
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			LogEmergencyError("failed to sync logs on shutdown (%s)", err)
		}
	}()

	shim.setLogger(logger)
	logger.Info("Logging initialized")

	health := NewHealth()
	processContainer := NewProcessContainer()

	serviceContainer := NewServiceContainer()
	_ = serviceContainer.Set("health", health)
	_ = serviceContainer.Set("logger", logger)
	_ = serviceContainer.Set("services", serviceContainer)
	_ = serviceContainer.Set("config", config)

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

	ctx := context.Background()
	if bs.contextFilter != nil {
		ctx = bs.contextFilter(ctx)
	}

	configs := loadConfig(processContainer, config, logger)

	if showHelp {
		description, err := describeConfiguration(config, configs, logger, &log.Config{})
		if err != nil {
			LogEmergencyError("failed to describe configuration (%s)", err)
			return 1
		}

		fmt.Println(description)
		return 0
	}

	if validateConfig(config, configs, logger) != nil {
		return 1
	}

	statusCode := 0
	for range runner.Run(ctx) {
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

func (bs *Bootstrapper) makeLogger(baseConfig Config, enable bool) (Logger, error) {
	if !enable {
		return NewNilLogger(), nil
	}

	logger, err := bs.loggingInitFunc(baseConfig)
	if err != nil {
		return nil, err
	}

	return logger.WithFields(bs.loggingFields), nil
}

func showHelp() bool {
	for _, arg := range os.Args[1:] {
		if arg == "--help" {
			return true
		}
	}

	return false
}
