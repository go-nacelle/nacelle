package nacelle

import (
	"github.com/efritz/nacelle/service"
)

type (
	ServiceContainer = service.Container

	// ServiceInitializerFunc is an InitializerFunc with a container argument.
	ServiceInitializerFunc func(config Config, container ServiceContainer) error
)

// WrapServiceInitializerFunc creates an InitializerFunc from a ServiceInitializerFunc and a container.
func WrapServiceInitializerFunc(container ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
	return InitializerFunc(func(config Config) error {
		return f(config, container)
	})
}
