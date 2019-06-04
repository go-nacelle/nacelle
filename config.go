package nacelle

import (
	"github.com/go-nacelle/config"
)

type (
	Config        = config.Config
	ConfigSourcer = config.Sourcer
	TagModifier   = config.TagModifier
)

var (
	NewConfig                   = config.NewConfig
	NewLoggingConfig            = config.NewLoggingConfig
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
