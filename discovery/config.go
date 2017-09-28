package discovery

import (
	"errors"
	"os"
	"time"
)

type (
	Config struct {
		DiscoveryAddr        string `env:"DISCOVERY_ADDR" required:"true"`
		DiscoveryBackend     string `env:"DISCOVERY_BACKEND" default:"consul"`
		RawDiscoveryTTL      int    `env:"DISCOVERY_TTL" default:"60"`
		RawDiscoveryInterval int    `env:"DISCOVERY_INTERVAL" default:"60"`
		DiscoveryPrefix      string `env:"DISCOVERY_PREFIX"`
		DiscoveryHost        string `env:"DISCOVERY_HOST"`
		DiscoveryPort        int    `env:"DISCOVERY_PORT"`

		DiscoveryTTL      time.Duration
		DiscoveryInterval time.Duration
	}

	configToken struct{}
)

var (
	ConfigToken       = configToken{}
	ErrIllegalBackend = errors.New("illegal discovery backend")
	ErrIllegalTTL     = errors.New("TTL must be greater than the refresh interval")
	ErrIllegalHost    = errors.New("hostname cannot be determined")
)

func (c *Config) PostLoad() error {
	if !isLegalBackend(c.DiscoveryBackend) {
		return ErrIllegalBackend
	}

	if c.RawDiscoveryTTL <= c.RawDiscoveryInterval {
		return ErrIllegalTTL
	}

	if c.DiscoveryHost == "" {
		c.DiscoveryHost = os.Getenv("HOST")

		if c.DiscoveryHost == "" {
			return ErrIllegalHost
		}
	}

	c.DiscoveryTTL = time.Duration(c.RawDiscoveryTTL) * time.Second
	c.DiscoveryInterval = time.Duration(c.RawDiscoveryInterval) * time.Second
	return nil
}

func isLegalBackend(backend string) bool {
	for _, whitelisted := range []string{"consul", "etcd", "zookeeper"} {
		if backend == whitelisted {
			return true
		}
	}

	return false
}
