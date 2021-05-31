package nacelle

import "github.com/go-nacelle/process/v2"

type (
	Health                = process.Health
	Initializer           = process.Initializer
	Finalizer             = process.Finalizer
	InitializerConfigFunc = process.InitializerConfigFunc
	InitializerFunc       = process.InitializerFunc
	Process               = process.Process
	ProcessConfigFunc     = process.ProcessConfigFunc
	ProcessContainer      = process.ProcessContainer
	RunnerConfigFunc      = process.RunnerConfigFunc
)

var (
	NewHealth                    = process.NewHealth
	NewProcessContainer          = process.NewProcessContainer
	WithHealthCheckInterval      = process.WithHealthCheckInterval
	WithInitializerPriority      = process.WithInitializerPriority
	WithInitializerContextFilter = process.WithInitializerContextFilter
	WithInitializerName          = process.WithInitializerName
	WithInitializerLogFields     = process.WithInitializerLogFields
	WithInitializerTimeout       = process.WithInitializerTimeout
	WithProcessPriority          = process.WithProcessPriority
	WithProcessInitTimeout       = process.WithProcessInitTimeout
	WithProcessContextFilter     = process.WithProcessContextFilter
	WithProcessName              = process.WithProcessName
	WithProcessLogFields         = process.WithProcessLogFields
	WithProcessShutdownTimeout   = process.WithProcessShutdownTimeout
	WithProcessStartTimeout      = process.WithProcessStartTimeout
	WithShutdownTimeout          = process.WithShutdownTimeout
	WithSilentExit               = process.WithSilentExit
	WithStartTimeout             = process.WithStartTimeout
)
