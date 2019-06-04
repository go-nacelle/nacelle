package service

import (
	"fmt"

	"github.com/go-nacelle/service"

	"github.com/go-nacelle/nacelle/logging"
)

type (
	// Container is a wrapper around services indexed by a unique
	// name. Services can be retrieved by name, or injected into a struct
	// by reading tagged fields.
	Container interface {
		// TODO - call it container
		service.ServiceContainer

		GetLogger() logging.Logger
	}

	serviceContainer struct {
		service.ServiceContainer
	}
)

// NewContainer creates an service container with itself registered
// to the key `container`.
func NewContainer() (Container, error) {
	container := &serviceContainer{
		ServiceContainer: service.NewServiceContainer(),
	}

	err := container.Set("container", container)
	return container, err
}

// GetLogger gets the service registered to the key `logger`. If no
// logger is registered, it will return an emergency logger instead.
// This method should always be safe to call.
func (c *serviceContainer) GetLogger() logging.Logger {
	if raw, err := c.Get("logger"); err == nil {
		return raw.(logging.Logger)
	}

	return logging.EmergencyLogger()
}

// Set registers a service with the given key. It is an error for
// a service to already be registered to this key. Additionally, it
// is an error to register an object that is not a Logger to the key
// `logger`.
func (c *serviceContainer) Set(key string, service interface{}) error {
	if _, ok := service.(logging.Logger); key == "logger" && !ok {
		return fmt.Errorf("logger instance is not a nacelle.Logger")
	}

	return c.ServiceContainer.Set(key, service)
}
