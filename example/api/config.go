package api

type (
	Config struct {
		RedisHost   string `env:"REDIS_HOST" default:"localhost"`
		RedisPort   int    `env:"REDIS_PORT" default:"6379"`
		RedisDB     int    `env:"REDIS_DB" default:"0"`
		RedisTTL    int    `env:"REDIS_TTL" default:"60"`
		RedisPrefix string `env:"REDIS_PREFIX" default:"secret"`
	}

	configToken struct{}
)

var ConfigToken = configToken{}
