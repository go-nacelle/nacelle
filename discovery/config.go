package discovery

import (
	"fmt"
	"os"
	"time"
)

type (
	Config struct {
		DiscoveryAddr        string `env:"DISCOVERY_ADDR" required:"true"`
		DiscoveryBackend     string `env:"DISCOVERY_BACKEND" default:"consul"`
		RawDiscoveryTTL      int    `env:"DISCOVERY_TTL" mask:"true" default:"120"`
		RawDiscoveryInterval int    `env:"DISCOVERY_INTERVAL" mask:"true" default:"60"`
		DiscoveryPrefix      string `env:"DISCOVERY_PREFIX"`
		DiscoveryHost        string `env:"DISCOVERY_HOST"`
		DiscoveryPort        int    `env:"DISCOVERY_PORT"`

		DiscoveryTTL      time.Duration
		DiscoveryInterval time.Duration
	}
)

var (
	ErrIllegalBackend = fmt.Errorf("illegal discovery backend")
	ErrIllegalTTL     = fmt.Errorf("TTL must be greater than the refresh interval")
	ErrIllegalHost    = fmt.Errorf("hostname cannot be determined")
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
