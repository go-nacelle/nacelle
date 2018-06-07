package secret

import (
	"fmt"

	"github.com/efritz/deepjoy"
	"github.com/efritz/nacelle"
)

func Init(config nacelle.Config, container *nacelle.ServiceContainer) error {
	secretConfig := &Config{}
	if err := config.Fetch(ConfigToken, secretConfig); err != nil {
		return err
	}

	logger := container.GetLogger()

	redisAddr := fmt.Sprintf(
		"%s:%d",
		secretConfig.RedisHost,
		secretConfig.RedisPort,
	)

	client := deepjoy.NewClient(
		redisAddr,
		deepjoy.WithDatabase(secretConfig.RedisDB),
		deepjoy.WithLogger(&logAdapter{logger}),
	)

	return container.Set("secret-service", &secretService{
		ttl:    secretConfig.RedisTTL,
		prefix: secretConfig.RedisPrefix,
		client: client,
	})
}
