package nacelle

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-nacelle/config/v3"
	"github.com/go-nacelle/log/v2"
	"github.com/go-nacelle/process/v2"
	"github.com/go-nacelle/service/v2"
)

// Bootstrapper wraps the entrypoint to the program.
type Bootstrapper struct {
	initFunc           AppInitFunc
	contextFilter      func(ctx context.Context) context.Context
	configSourcer      ConfigSourcer
	configMaskedKeys   []string
	loggingInitFunc    LoggingInitFunc
	loggingFields      LogFields
	machineConfigFuncs []MachineConfigFunc
}

type bootstrapperConfig struct {
	contextFilter      func(ctx context.Context) context.Context
	configSourcer      ConfigSourcer
	configMaskedKeys   []string
	loggingInitFunc    LoggingInitFunc
	loggingFields      LogFields
	machineConfigFuncs []MachineConfigFunc
}

// AppInitFunc is an program entrypoint called after performing initial
// configuration loading, sanity checks, and setting up loggers. This
// function should register initializers and processes and inject values
// into the service container where necessary.
type AppInitFunc func(context.Context, *ProcessContainerBuilder, *ServiceContainer) error

// ServiceInitializerFunc is an InitializerFunc with a service container argument.
type ServiceInitializerFunc func(ctx context.Context, container *ServiceContainer) error

// WrapServiceInitializerFunc creates an InitializerFunc from a ServiceInitializerFunc and a container.
func WrapServiceInitializerFunc(container *ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
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
		initFunc:           initFunc,
		contextFilter:      config.contextFilter,
		configSourcer:      config.configSourcer,
		configMaskedKeys:   config.configMaskedKeys,
		loggingInitFunc:    config.loggingInitFunc,
		loggingFields:      config.loggingFields,
		machineConfigFuncs: config.machineConfigFuncs,
	}
}

// Boot will initialize services and return a status code. This
// method does not return in any meaningful way (it blocks until
// the associated process runner has completed).
func (bs *Bootstrapper) Boot() int {
	showHelp := showHelp()

	shim := &logShim{}
	cfg := NewConfig(
		bs.configSourcer,
		config.WithLogger(shim),
		config.WithMaskedKeys(bs.configMaskedKeys),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cfg.Init(); err != nil {
		LogEmergencyError("failed to initialize config (%s)", err)
		return 1
	}

	logger, err := bs.makeLogger(cfg, !showHelp)
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

	// Add default services to service container
	serviceContainer := NewServiceContainer()
	_ = serviceContainer.Set("health", health)
	_ = serviceContainer.Set("logger", logger)
	_ = serviceContainer.Set("services", serviceContainer)
	_ = serviceContainer.Set("config", cfg)

	// Add default services to context
	ctx = process.ContextWithHealth(ctx, health)
	ctx = log.WithLogger(ctx, logger)
	ctx = service.WithContainer(ctx, serviceContainer)
	ctx = config.WithConfig(ctx, cfg)

	// Run the context filter after we've added the services to the context in
	// case the user wants to tweak those services.
	if bs.contextFilter != nil {
		ctx = bs.contextFilter(ctx)
	}

	processContainerBuilder := process.NewContainerBuilder()
	if err := bs.initFunc(ctx, processContainerBuilder, serviceContainer); err != nil {
		logger.Error("Failed to run initialization function (%s)", err.Error())
		return 1
	}
	processContainer := processContainerBuilder.Build(process.WithMetaLogger(&logAdapter{logger}))

	configs := loadConfig(processContainer, cfg, logger)

	if showHelp {
		description, err := describeConfiguration(cfg, configs, logger, &log.Config{})
		if err != nil {
			LogEmergencyError("failed to describe configuration (%s)", err)
			return 1
		}

		fmt.Println(description)
		return 0
	}

	if validateConfig(cfg, configs, logger) != nil {
		return 1
	}

	defaultConfigs := []process.MachineConfigFunc{
		process.WithHealth(health),
		process.WithInjecter(newInjectHook(serviceContainer, logger)),
	}
	state := process.Run(ctx, processContainer, append(defaultConfigs, bs.machineConfigFuncs...)...)

	ch := make(chan os.Signal, 2)
	for _, s := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		signal.Notify(ch, s)
	}

	go func() {
		<-ch
		logger.Info("Received signal")
		state.Shutdown(context.Background())

		<-ch
		logger.Error("Received second signal, no longer waiting for graceful exit")
		cancel()

		<-ch
		os.Exit(1)
	}()

	statusCode := 0
	if !state.Wait(ctx) {
		statusCode = 1
		for _, err := range state.Errors() {
			logger.Error("%s", err)
		}
	}

	logger.Info("All processes have stopped")
	return statusCode
}

// BootAndExit calls Boot and sets the program return code on halt. This
// method does not return.
func (bs *Bootstrapper) BootAndExit() {
	os.Exit(bs.Boot())
}

func (bs *Bootstrapper) makeLogger(baseConfig *Config, enable bool) (Logger, error) {
	if !enable {
		return NewNilLogger(), nil
	}

	logger, err := bs.loggingInitFunc(baseConfig)
	if err != nil {
		return nil, err
	}

	return logger.WithFields(bs.loggingFields), nil
}

func newInjectHook(serviceContainer *ServiceContainer, logger Logger) process.Injecter {
	return process.InjecterFunc(func(ctx context.Context, meta *process.Meta) error {
		serviceContainer, err := replaceLoggerService(serviceContainer, logger.WithFields(LogFields(meta.Metadata())))
		if err != nil {
			return err
		}

		return service.Inject(ctx, serviceContainer, meta.Wrapped())
	})
}

func replaceLoggerService(serviceContainer *ServiceContainer, logger Logger) (*ServiceContainer, error) {
	// Update logger instance
	serviceContainer, err := overlay(serviceContainer, "logger", logger)
	if err != nil {
		return nil, err
	}

	// Update self-reference
	serviceContainer, err = overlay(serviceContainer, "services", serviceContainer)
	if err != nil {
		return nil, err
	}

	return serviceContainer, nil
}

func overlay(serviceContainer *ServiceContainer, key, service interface{}) (*ServiceContainer, error) {
	return serviceContainer.WithValues(map[interface{}]interface{}{key: service})
}

func showHelp() bool {
	for _, arg := range os.Args[1:] {
		if arg == "--help" {
			return true
		}
	}

	return false
}

type logAdapter struct {
	log.Logger
}

func (adapter *logAdapter) WithFields(fields process.LogFields) process.Logger {
	return &logAdapter{adapter.Logger.WithFields(LogFields(fields))}
}
