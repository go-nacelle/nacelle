package nacelle

import (
	"github.com/efritz/nacelle/config"
	"github.com/efritz/zubrin"
)

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
	NewGlobSourcer              = zubrin.NewGlobSourcer
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
