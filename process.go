package nacelle

import "github.com/go-nacelle/nacelle/process"

type (
	Process               = process.Process
	Initializer           = process.Initializer
	InitializerFunc       = process.InitializerFunc
	ParallelInitializer   = process.ParallelInitializer
	ProcessContainer      = process.Container
	Health                = process.Health
	ProcessConfigFunc     = process.ProcessConfigFunc
	InitializerConfigFunc = process.InitializerConfigFunc
	RunnerConfigFunc      = process.RunnerConfigFunc

	// ServiceInitializerFunc is an InitializerFunc with a service container argument.
	ServiceInitializerFunc func(config Config, container ServiceContainer) error
)

var (
	NewParallelInitializer     = process.NewParallelInitializer
	WithStartTimeout           = process.WithStartTimeout
	WithHealthCheckBackoff     = process.WithHealthCheckBackoff
	WithShutdownTimeout        = process.WithShutdownTimeout
	WithInitializerName        = process.WithInitializerName
	WithProcessName            = process.WithProcessName
	WithPriority               = process.WithPriority
	WithSilentExit             = process.WithSilentExit
	WithInitializerTimeout     = process.WithInitializerTimeout
	WithProcessInitTimeout     = process.WithProcessInitTimeout
	WithProcessStartTimeout    = process.WithProcessStartTimeout
	WithProcessShutdownTimeout = process.WithProcessShutdownTimeout
)

// WrapServiceInitializerFunc creates an InitializerFunc from a ServiceInitializerFunc and a container.
func WrapServiceInitializerFunc(container ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
	return InitializerFunc(func(config Config) error {
		return f(config, container)
	})
}
