package nacelle

import (
	"fmt"

	"github.com/efritz/bussard"
)

type (
	// ServiceContainer is a wrapper around services indexed by a unique
	// name. Services can be retrieved by name, or injected into a struct
	// by reading tagged fields.
	ServiceContainer struct {
		bussard.ServiceContainer
	}

	// ServiceInitializerFunc is an InitializerFunc with a container argument.
	ServiceInitializerFunc func(config Config, container *ServiceContainer) error
)

// WrapServiceInitializerFunc creates an InitializerFunc from a ServiceInitializerFunc and a container.
func WrapServiceInitializerFunc(container *ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
	return InitializerFunc(func(config Config) error {
		return f(config, container)
	})
}

// MakeServiceContainer creates an service container with itself
// regsitered to the key `container`.
func MakeServiceContainer() (*ServiceContainer, error) {
	container := &ServiceContainer{bussard.NewServiceContainer()}
	err := container.Set("container", container)
	return container, err
}

// GetLogger gets the service registered to the key `logger`. If no
// logger is registered, it will return an emergency logger instead.
// This method should always be safe to call.
func (c *ServiceContainer) GetLogger() Logger {
	if raw, err := c.Get("logger"); err == nil {
		return raw.(Logger)
	}

	return emergencyLogger()
}

// Set registers a service with the given key. It is an error for
// a service to already be registered to this key. Additionally, it
// is an error to register an object that is not a Logger to the key
// `logger`.
func (c *ServiceContainer) Set(key string, service interface{}) error {
	if _, ok := service.(Logger); key == "logger" && !ok {
		return fmt.Errorf("logger instance is not a nacelle.Logger")
	}

	return c.ServiceContainer.Set(key, service)
}
