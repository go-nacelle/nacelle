package main

import (
	"net/http"

	"github.com/efritz/nacelle"
	basehttp "github.com/efritz/nacelle/base/http"
)

func setupServer(config nacelle.Config, server *http.Server) error {
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!\n"))
	})

	return nil
}

//
//

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	processes.RegisterProcess(basehttp.NewServer(basehttp.ServerInitializerFunc(setupServer)))
	return nil
}

//
//

func main() {
	nacelle.NewBootstrapper("http-example", setup).BootAndExit()
}
