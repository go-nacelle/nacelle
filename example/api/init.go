package api

import (
	"fmt"

	"github.com/efritz/deepjoy"
	"github.com/efritz/nacelle"
)

func Init(config nacelle.Config, container *nacelle.ServiceContainer) error {
	cfg, err := config.Get(ConfigToken)
	if err != nil {
		return err
	}

	apiConfig := cfg.(*Config)

	logger, err := container.Get("logger")
	if err != nil {
		return err
	}

	redisAddr := fmt.Sprintf(
		"%s:%d",
		apiConfig.RedisHost,
		apiConfig.RedisPort,
	)

	client := deepjoy.NewClient(
		redisAddr,
		deepjoy.WithDatabase(apiConfig.RedisDB),
		deepjoy.WithLogger(&logAdapter{
			logger: logger.(nacelle.Logger),
		}),
	)

	service := &secretService{
		ttl:    apiConfig.RedisTTL,
		prefix: apiConfig.RedisPrefix,
		client: client,
	}

	return container.Set("secret_service", service)
}
