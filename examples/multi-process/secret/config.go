package secret

type Config struct {
	RedisHost   string `env:"redis_host" default:"localhost"`
	RedisPort   int    `env:"redis_port" default:"6379"`
	RedisDB     int    `env:"redis_db" default:"0"`
	RedisTTL    int    `env:"redis_ttl" default:"60"`
	RedisPrefix string `env:"redis_prefix" default:"secret"`
}
