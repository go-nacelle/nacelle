package nacelle

import (
	"github.com/go-nacelle/config"
	"github.com/go-nacelle/log"
	"github.com/go-nacelle/process"
	"github.com/go-nacelle/service"
)

// TODO - ensure ServiceContainer instances are named consistently

type (
	Config                = config.Config
	ConfigSourcer         = config.Sourcer
	Health                = process.Health
	Initializer           = process.Initializer
	InitializerConfigFunc = process.InitializerConfigFunc
	InitializerFunc       = process.InitializerFunc
	LogFields             = log.Fields
	Logger                = log.Logger
	LogLevel              = log.LogLevel
	ParallelInitializer   = process.ParallelInitializer
	Process               = process.Process
	ProcessConfigFunc     = process.ProcessConfigFunc
	ProcessContainer      = process.ProcessContainer
	RunnerConfigFunc      = process.RunnerConfigFunc
	ServiceContainer      = service.ServiceContainer
	TagModifier           = config.TagModifier
)

const (
	LevelDebug   = log.LevelDebug
	LevelError   = log.LevelError
	LevelFatal   = log.LevelFatal
	LevelInfo    = log.LevelInfo
	LevelWarning = log.LevelWarning
)

var (
	EmergencyLogger             = log.EmergencyLogger
	LogEmergencyError           = log.LogEmergencyError
	LogEmergencyErrors          = log.LogEmergencyErrors
	NewConfig                   = config.NewConfig
	NewDefaultTagSetter         = config.NewDefaultTagSetter
	NewDirectorySourcer         = config.NewDirectorySourcer
	NewEnvSourcer               = config.NewEnvSourcer
	NewEnvTagPrefixer           = config.NewEnvTagPrefixer
	NewFileSourcer              = config.NewFileSourcer
	NewFileTagPrefixer          = config.NewFileTagPrefixer
	NewGlobSourcer              = config.NewGlobSourcer
	NewHealth                   = process.NewHealth
	NewMultiSourcer             = config.NewMultiSourcer
	NewNilLogger                = log.NewNilLogger
	NewOptionalDirectorySourcer = config.NewOptionalDirectorySourcer
	NewOptionalFileSourcer      = config.NewOptionalFileSourcer
	NewParallelInitializer      = process.NewParallelInitializer
	NewProcessContainer         = process.NewProcessContainer
	NewReplayAdapter            = log.NewReplayAdapter
	NewRollupAdapter            = log.NewRollupAdapter
	NewServiceContainer         = service.NewServiceContainer
	NewTOMLFileSourcer          = config.NewTOMLFileSourcer
	NewYAMLFileSourcer          = config.NewYAMLFileSourcer
	WithHealthCheckBackoff      = process.WithHealthCheckBackoff
	WithInitializerName         = process.WithInitializerName
	WithInitializerTimeout      = process.WithInitializerTimeout
	WithPriority                = process.WithPriority
	WithProcessInitTimeout      = process.WithProcessInitTimeout
	WithProcessName             = process.WithProcessName
	WithProcessShutdownTimeout  = process.WithProcessShutdownTimeout
	WithProcessStartTimeout     = process.WithProcessStartTimeout
	WithShutdownTimeout         = process.WithShutdownTimeout
	WithSilentExit              = process.WithSilentExit
	WithStartTimeout            = process.WithStartTimeout
)
