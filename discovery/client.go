package discovery

import (
	"github.com/efritz/reception"

	"github.com/efritz/nacelle"
)

type initFunc func(*Config, nacelle.Logger) (reception.Client, error)

var (
	initializers = map[string]initFunc{
		"consul":    initConsul,
		"etd":       initEtcd,
		"zookeeper": initZk,
	}
)

func makeClient(config nacelle.Config, container nacelle.ServiceContainer) (reception.Client, error) {
	discoveryConfig := &Config{}
	if err := config.Load(discoveryConfig); err != nil {
		return nil, err
	}

	return initializers[discoveryConfig.DiscoveryBackend](
		discoveryConfig,
		container.GetLogger(),
	)
}

func initConsul(config *Config, logger nacelle.Logger) (reception.Client, error) {
	return reception.DialConsul(
		config.DiscoveryAddr,
		reception.WithHost(config.DiscoveryHost),
		reception.WithPort(config.DiscoveryPort),
		reception.WithCheckTimeout(config.DiscoveryTTL),
		reception.WithCheckInterval(config.DiscoveryInterval),
		reception.WithCheckDeregisterTimeout(config.DiscoveryInterval),
		reception.WithLogger(&logAdapter{logger}),
	)
}

func initEtcd(config *Config, logger nacelle.Logger) (reception.Client, error) {
	return reception.DialEtcd(
		config.DiscoveryAddr,
		reception.WithEtcdPrefix(config.DiscoveryPrefix),
		reception.WithTTL(config.DiscoveryTTL),
		reception.WithRefreshInterval(config.DiscoveryInterval),
	)
}

func initZk(config *Config, logger nacelle.Logger) (reception.Client, error) {
	return reception.DialExhibitor(
		config.DiscoveryAddr,
		reception.WithZkPrefix(config.DiscoveryPrefix),
	)
}
