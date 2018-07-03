package main

//go:generate protoc grpc/secret.proto --go_out=plugins=grpc:grpc -I grpc/

import (
	"github.com/efritz/nacelle"
	basegrpc "github.com/efritz/nacelle/base/grpc"
	basehttp "github.com/efritz/nacelle/base/http"

	"github.com/efritz/nacelle/examples/multi-process/grpc"
	"github.com/efritz/nacelle/examples/multi-process/http"
	"github.com/efritz/nacelle/examples/multi-process/secret"
)

func setupConfigs(config nacelle.Config) error {
	config.MustRegister(secret.ConfigToken, &secret.Config{})
	config.MustRegister(basehttp.ConfigToken, &basehttp.Config{})
	config.MustRegister(basegrpc.ConfigToken, &basegrpc.Config{})
	return nil
}

func setupProcesses(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	initSecret := nacelle.WrapServiceInitializerFunc(services, secret.Init)

	processes.RegisterInitializer(initSecret, nacelle.WithInitializerName("secret"))
	processes.RegisterProcess(basehttp.NewServer(http.NewEndpointSet()), nacelle.WithProcessName("http"))
	processes.RegisterProcess(basegrpc.NewServer(grpc.NewEndpointSet()), nacelle.WithProcessName("grpc"))
	return nil
}

func main() {
	boostrapper := nacelle.NewBootstrapper(
		"multi-process-example",
		setupConfigs,
		setupProcesses,
		nacelle.WithLoggingFields(nacelle.Fields{
			"app_name": "multi-process-example",
		}),
	)

	boostrapper.BootAndExit()
}
