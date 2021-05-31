package nacelle

import "github.com/go-nacelle/process/v2"

type (
	Health                  = process.Health
	HealthComponentStatus   = process.HealthComponentStatus
	Injecter                = process.Injecter
	Initializer             = process.Initializer
	Runner                  = process.Runner
	Stopper                 = process.Stopper
	Finalizer               = process.Finalizer
	InitializerFunc         = process.InitializerFunc
	InjecterFunc            = process.InjecterFunc
	RunnerFunc              = process.RunnerFunc
	ProcessContainer        = process.Container
	ProcessContainerBuilder = process.ContainerBuilder
	MachineConfigFunc       = process.MachineConfigFunc
	MetaConfigFunc          = process.MetaConfigFunc
)

var (
	NewHealth               = process.NewHealth
	WithInjecter            = process.WithInjecter
	WithHealth              = process.WithHealth
	WithMetaHealth          = process.WithMetaHealth
	WithMetaHealthKey       = process.WithMetaHealthKey
	WithMetaContext         = process.WithMetaContext
	WithMetaName            = process.WithMetaName
	WithMetaPriority        = process.WithMetaPriority
	WithMetadata            = process.WithMetadata
	WithEarlyExit           = process.WithEarlyExit
	WithMetaInitTimeout     = process.WithMetaInitTimeout
	WithMetaStartupTimeout  = process.WithMetaStartupTimeout
	WithMetaStopTimeout     = process.WithMetaStopTimeout
	WithMetaShutdownTimeout = process.WithMetaShutdownTimeout
	WithMetaFinalizeTimeout = process.WithMetaFinalizeTimeout
	WithMetaLogger          = process.WithMetaLogger
)
