package nacelle

import (
	"github.com/go-nacelle/process"
)

type (
	Health                = process.Health
	Initializer           = process.Initializer
	InitializerConfigFunc = process.InitializerConfigFunc
	InitializerFunc       = process.InitializerFunc
	ParallelInitializer   = process.ParallelInitializer
	Process               = process.Process
	ProcessConfigFunc     = process.ProcessConfigFunc
	ProcessContainer      = process.ProcessContainer
	RunnerConfigFunc      = process.RunnerConfigFunc
)

var (
	NewHealth                  = process.NewHealth
	NewParallelInitializer     = process.NewParallelInitializer
	NewProcessContainer        = process.NewProcessContainer
	WithHealthCheckBackoff     = process.WithHealthCheckBackoff
	WithInitializerName        = process.WithInitializerName
	WithInitializerTimeout     = process.WithInitializerTimeout
	WithPriority               = process.WithPriority
	WithProcessInitTimeout     = process.WithProcessInitTimeout
	WithProcessName            = process.WithProcessName
	WithProcessShutdownTimeout = process.WithProcessShutdownTimeout
	WithProcessStartTimeout    = process.WithProcessStartTimeout
	WithShutdownTimeout        = process.WithShutdownTimeout
	WithSilentExit             = process.WithSilentExit
	WithStartTimeout           = process.WithStartTimeout
)
