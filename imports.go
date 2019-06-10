package nacelle

import (
	"github.com/go-nacelle/config"
	"github.com/go-nacelle/log"
	"github.com/go-nacelle/process"
	"github.com/go-nacelle/service"
)

type (
	Process               = process.Process
	Initializer           = process.Initializer
	InitializerFunc       = process.InitializerFunc
	ParallelInitializer   = process.ParallelInitializer
	ProcessContainer      = process.ProcessContainer
	Health                = process.Health
	ProcessConfigFunc     = process.ProcessConfigFunc
	InitializerConfigFunc = process.InitializerConfigFunc
	RunnerConfigFunc      = process.RunnerConfigFunc
)

var NewProcessContainer = process.NewProcessContainer

var (
	NewParallelInitializer     = process.NewParallelInitializer
	WithStartTimeout           = process.WithStartTimeout
	WithHealthCheckBackoff     = process.WithHealthCheckBackoff
	NewHealth                  = process.NewHealth
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

// TODO - call it services everywhere
type ServiceContainer = service.ServiceContainer

var NewServiceContainer = service.NewServiceContainer

type (
	Logger    = log.Logger
	LogLevel  = log.LogLevel
	LogFields = log.Fields
)

const (
	LevelFatal   = log.LevelFatal
	LevelError   = log.LevelError
	LevelWarning = log.LevelWarning
	LevelInfo    = log.LevelInfo
	LevelDebug   = log.LevelDebug
)

var (
	NewNilLogger       = log.NewNilLogger
	NewReplayAdapter   = log.NewReplayAdapter
	NewRollupAdapter   = log.NewRollupAdapter
	LogEmergencyError  = log.LogEmergencyError
	LogEmergencyErrors = log.LogEmergencyErrors
	EmergencyLogger    = log.EmergencyLogger
)

type (
	Config        = config.Config
	ConfigSourcer = config.Sourcer
	TagModifier   = config.TagModifier
)

var (
	NewConfig                   = config.NewConfig
	NewEnvSourcer               = config.NewEnvSourcer
	NewFileSourcer              = config.NewFileSourcer
	NewGlobSourcer              = config.NewGlobSourcer
	NewOptionalFileSourcer      = config.NewOptionalFileSourcer
	NewDirectorySourcer         = config.NewDirectorySourcer
	NewOptionalDirectorySourcer = config.NewOptionalDirectorySourcer
	NewYAMLFileSourcer          = config.NewYAMLFileSourcer
	NewTOMLFileSourcer          = config.NewTOMLFileSourcer
	NewMultiSourcer             = config.NewMultiSourcer
	NewEnvTagPrefixer           = config.NewEnvTagPrefixer
	NewFileTagPrefixer          = config.NewFileTagPrefixer
	NewDefaultTagSetter         = config.NewDefaultTagSetter
)
