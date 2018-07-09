package nacelle

import "github.com/efritz/nacelle/process"

type (
	Process               = process.Process
	Initializer           = process.Initializer
	InitializerFunc       = process.InitializerFunc
	ProcessContainer      = process.Container
	Health                = process.Health
	ProcessConfigFunc     = process.ProcessConfigFunc
	InitializerConfigFunc = process.InitializerConfigFunc
	RunnerConfigFunc      = process.RunnerConfigFunc

	// ServiceInitializerFunc is an InitializerFunc with a service container argument.
	ServiceInitializerFunc func(config Config, container ServiceContainer) error
)

var (
	WithShutdownTimeout    = process.WithShutdownTimeout
	WithInitializerName    = process.WithInitializerName
	WithProcessName        = process.WithProcessName
	WithPriority           = process.WithPriority
	WithSilentExit         = process.WithSilentExit
	WithInitializerTimeout = process.WithInitializerTimeout
	WithProcessInitTimeout = process.WithProcessInitTimeout
)

// WrapServiceInitializerFunc creates an InitializerFunc from a ServiceInitializerFunc and a container.
func WrapServiceInitializerFunc(container ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
	return InitializerFunc(func(config Config) error {
		return f(config, container)
	})
}
