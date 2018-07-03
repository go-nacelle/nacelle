package secret

import (
	"fmt"

	"github.com/efritz/deepjoy"
	"github.com/efritz/nacelle"
)

func Init(config nacelle.Config, services nacelle.ServiceContainer) error {
	secretConfig := &Config{}
	if err := config.Fetch(ConfigToken, secretConfig); err != nil {
		return err
	}

	var (
		host   = secretConfig.RedisHost
		port   = secretConfig.RedisPort
		db     = secretConfig.RedisDB
		ttl    = secretConfig.RedisTTL
		prefix = secretConfig.RedisPrefix

		secrets = &secretService{
			ttl:    ttl,
			prefix: prefix,
			client: deepjoy.NewClient(
				fmt.Sprintf("%s:%d", host, port),
				deepjoy.WithDatabase(db),
				deepjoy.WithLogger(&logAdapter{Logger: services.GetLogger()}),
			),
		}
	)

	return services.Set("secret-service", secrets)
}
