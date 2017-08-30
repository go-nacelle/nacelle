package main

import (
	"github.com/efritz/nacelle"

	"github.com/efritz/nacelle/example/api"
	"github.com/efritz/nacelle/example/grpc"
	"github.com/efritz/nacelle/example/http"
)

func main() {
	configs := map[interface{}]interface{}{
		api.ConfigToken:  &api.Config{},
		http.ConfigToken: &http.Config{},
		grpc.ConfigToken: &grpc.Config{},
	}

	nacelle.Boot("app", configs, setup)
}

func setup(runner *nacelle.ProcessRunner, container *nacelle.ServiceContainer) error {
	runner.RegisterInitializer(nacelle.WrapServiceInitializerFunc(api.Init, container))
	runner.RegisterProcess(http.NewServer(), 1)
	runner.RegisterProcess(grpc.NewServer(), 1)
	return nil
}
