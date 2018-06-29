package main

import (
	"github.com/efritz/nacelle"
	basegrpc "github.com/efritz/nacelle/base/grpc"
	basehttp "github.com/efritz/nacelle/base/http"

	"github.com/efritz/nacelle/examples/multi-service/grpc"
	"github.com/efritz/nacelle/examples/multi-service/http"
	"github.com/efritz/nacelle/examples/multi-service/secret"
)

func setupConfigs(config nacelle.Config) error {
	config.MustRegister(secret.ConfigToken, &secret.Config{})
	config.MustRegister(basehttp.ConfigToken, &basehttp.Config{})
	config.MustRegister(basegrpc.ConfigToken, &basegrpc.Config{})
	return nil
}

func setupProcesses(runner nacelle.ProcessContainer, container nacelle.ServiceContainer) error {
	runner.RegisterInitializer(nacelle.WrapServiceInitializerFunc(container, secret.Init))
	runner.RegisterProcess(basehttp.NewServer(http.NewEndpointSet()), nacelle.WithProcessName("http"))
	runner.RegisterProcess(basegrpc.NewServer(grpc.NewEndpointSet()), nacelle.WithProcessName("grpc"))
	return nil
}

func main() {
	boostrapper := nacelle.NewBootstrapper(
		"app",
		setupConfigs,
		setupProcesses,
		nacelle.WithLoggingFields(nacelle.Fields{
			"app_name": "app",
		}),
	)

	boostrapper.BootAndExit()
}
