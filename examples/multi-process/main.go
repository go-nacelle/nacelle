package main

//go:generate protoc grpc/secret.proto --go_out=plugins=grpc:grpc -I grpc/

import (
	"github.com/go-nacelle/nacelle"
	basegrpc "github.com/go-nacelle/nacelle/base/grpc"
	basehttp "github.com/go-nacelle/nacelle/base/http"

	"github.com/go-nacelle/nacelle/examples/multi-process/grpc"
	"github.com/go-nacelle/nacelle/examples/multi-process/http"
	"github.com/go-nacelle/nacelle/examples/multi-process/secret"
)

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	initSecret := nacelle.WrapServiceInitializerFunc(services, secret.Init)

	processes.RegisterInitializer(initSecret, nacelle.WithInitializerName("secret"))
	processes.RegisterProcess(basehttp.NewServer(http.NewEndpointSet()), nacelle.WithProcessName("http"))
	processes.RegisterProcess(basegrpc.NewServer(grpc.NewEndpointSet()), nacelle.WithProcessName("grpc"))
	return nil
}

func main() {
	boostrapper := nacelle.NewBootstrapper(
		"multi-process-example",
		setup,
		nacelle.WithLoggingFields(nacelle.LogFields{
			"app_name": "multi-process-example",
		}),
	)

	boostrapper.BootAndExit()
}
