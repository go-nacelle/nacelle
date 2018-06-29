package http

import (
	"net"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/service"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ConfigSuite{})
		s.AddSuite(&ServerSuite{})
	})
}

//
// Config

type emptyConfig struct{}

func makeConfig(token, base interface{}) nacelle.Config {
	config := nacelle.NewEnvConfig("")
	config.Register(token, base)
	config.Load()

	return config
}

//
//  Injection

type A struct{ X int }
type B struct{ X float64 }

func makeBadContainer() nacelle.ServiceContainer {
	container, _ := service.NewContainer()
	container.Set("A", &B{})
	return container
}

//
// Server Helpers

func getDynamicPort(listener net.Listener) int {
	return listener.Addr().(*net.TCPAddr).Port
}
