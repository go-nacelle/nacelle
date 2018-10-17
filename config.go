package nacelle

import "github.com/efritz/nacelle/config"

type (
	Config        = config.Config
	LoggingConfig = config.LoggingConfig
	ConfigSourcer = config.Sourcer
	TagModifier   = config.TagModifier
)

var (
	NewConfig                   = config.NewConfig
	NewLoggingConfig            = config.NewLoggingConfig
	NewEnvSourcer               = config.NewEnvSourcer
	NewFileSourcer              = config.NewFileSourcer
	NewOptionalFileSourcer      = config.NewOptionalFileSourcer
	NewDirectorySourcer         = config.NewDirectorySourcer
	NewOptionalDirectorySourcer = config.NewOptionalDirectorySourcer
	NewYAMLFileSourcer          = config.NewYAMLFileSourcer
	NewTOMLFileSourcer          = config.NewTOMLFileSourcer
	NewMultiSourcer             = config.NewMultiSourcer
	NewEnvTagPrefixer           = config.NewEnvTagPrefixer
	NewDefaultTagSetter         = config.NewDefaultTagSetter
)
