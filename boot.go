package nacelle

import "os"

type AppInitFunc func(*ProcessRunner, *ServiceContainer) error

func Boot(name string, configs map[interface{}]interface{}, initFunc AppInitFunc) {
	os.Exit(boot(name, configs, initFunc))
}

func boot(name string, configs map[interface{}]interface{}, initFunc AppInitFunc) int {
	var (
		container = NewServiceContainer()
		runner    = NewProcessRunner(container)
		config    = NewEnvConfig(name)
	)

	if configs == nil {
		configs = map[interface{}]interface{}{}
	}

	for key, obj := range configs {
		if err := config.Register(key, obj); err != nil {
			emergencyLogger().Error(nil, err.Error())
			return 1
		}
	}

	if err := config.Register(LoggingConfigToken, &LoggingConfig{}); err != nil {
		emergencyLogger().Error(nil, err.Error())
		return 1
	}

	if errs := config.Load(); len(errs) > 0 {
		logger := emergencyLogger()

		for _, err := range errs {
			logger.Error(nil, err.Error())
		}

		return 1
	}

	logger, err := InitLogging(config)
	if err != nil {
		emergencyLogger().Error(nil, err.Error())
		return 1
	}

	defer logger.Sync()
	logger.Info(nil, "Logging initialized")

	if err := container.Set("logger", logger); err != nil {
		logger.Error(nil, err.Error())
		return 1
	}

	m, err := config.ToMap()
	if err != nil {
		logger.Error(nil, err.Error())
		return 1
	}

	logger.Info(m, "Process starting")

	if err := initFunc(runner, container); err != nil {
		logger.Error(nil, err.Error())
		return 1
	}

	statusCode := 0
	for err := range runner.Run(config, logger) {
		statusCode = 1
		logger.Error(nil, err.Error())
	}

	logger.Info(nil, "All processes have stopped")
	return statusCode
}
