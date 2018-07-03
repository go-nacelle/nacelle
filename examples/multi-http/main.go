package main

import (
	"net/http"

	"github.com/efritz/nacelle"
	basehttp "github.com/efritz/nacelle/base/http"
)

func setupServerA(config nacelle.Config, server *http.Server) error {
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server A!\n"))
	})

	return nil
}

func setupServerB(config nacelle.Config, server *http.Server) error {
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server B!\n"))
	})

	return nil
}

//
//

var (
	tokenA = basehttp.NewConfigToken("A")
	tokenB = basehttp.NewConfigToken("B")
)

func setupConfigs(config nacelle.Config) error {
	modifiedConfig := nacelle.MustApplyTagModifiers(
		&basehttp.Config{},
		nacelle.NewEnvTagPrefixer("b"),
		nacelle.NewDefaultTagSetter("HTTPPort", "5001"),
	)

	config.MustRegister(tokenA, &basehttp.Config{})
	config.MustRegister(tokenB, modifiedConfig)
	return nil
}

func setupProcesses(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	serverA := basehttp.NewServer(basehttp.ServerInitializerFunc(setupServerA), basehttp.WithConfigToken(tokenA))
	serverB := basehttp.NewServer(basehttp.ServerInitializerFunc(setupServerB), basehttp.WithConfigToken(tokenB))

	processes.RegisterProcess(serverA, nacelle.WithProcessName("http-a"))
	processes.RegisterProcess(serverB, nacelle.WithProcessName("http-b"))
	return nil
}

//
//

func main() {
	nacelle.NewBootstrapper("multi-http-example", setupConfigs, setupProcesses).BootAndExit()
}
