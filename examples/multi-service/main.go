package main

import (
	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/process"

	"github.com/efritz/nacelle/examples/multi-service/grpc"
	"github.com/efritz/nacelle/examples/multi-service/http"
	"github.com/efritz/nacelle/examples/multi-service/secret"
)

func setupConfigs(config nacelle.Config) error {
	config.MustRegister(secret.ConfigToken, &secret.Config{})
	config.MustRegister(process.HTTPConfigToken, &process.HTTPConfig{})
	config.MustRegister(process.GRPCConfigToken, &process.GRPCConfig{})
	return nil
}

func setupProcesses(runner *nacelle.ProcessRunner, container nacelle.ServiceContainer) error {
	runner.RegisterInitializer(nacelle.WrapServiceInitializerFunc(container, secret.Init))
	runner.RegisterProcess(process.NewHTTPServer(http.NewEndpointSet()), nacelle.WithProcessName("http"))
	runner.RegisterProcess(process.NewGRPCServer(grpc.NewEndpointSet()), nacelle.WithProcessName("grpc"))
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
