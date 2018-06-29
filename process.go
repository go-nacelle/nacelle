package nacelle

import "github.com/efritz/nacelle/process"

type (
	Process               = process.Process
	Initializer           = process.Initializer
	InitializerFunc       = process.InitializerFunc
	ProcessContainer      = process.Container
	ProcessConfigFunc     = process.ProcessConfigFunc
	InitializerConfigFunc = process.InitializerConfigFunc
	RunnerConfigFunc      = process.RunnerConfigFunc
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
