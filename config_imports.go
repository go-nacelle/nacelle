package nacelle

import (
	"github.com/go-nacelle/config"
)

type (
	Config        = config.Config
	ConfigSourcer = config.Sourcer
	FileParser    = config.FileParser
	TagModifier   = config.TagModifier
)

var (
	NewConfig                   = config.NewConfig
	NewDefaultTagSetter         = config.NewDefaultTagSetter
	NewDirectorySourcer         = config.NewDirectorySourcer
	NewEnvSourcer               = config.NewEnvSourcer
	NewEnvTagPrefixer           = config.NewEnvTagPrefixer
	NewFileSourcer              = config.NewFileSourcer
	NewFileTagPrefixer          = config.NewFileTagPrefixer
	NewFileTagSetter            = config.NewFileTagSetter
	NewGlobSourcer              = config.NewGlobSourcer
	NewMultiSourcer             = config.NewMultiSourcer
	NewOptionalDirectorySourcer = config.NewOptionalDirectorySourcer
	NewOptionalFileSourcer      = config.NewOptionalFileSourcer
	NewTOMLFileSourcer          = config.NewTOMLFileSourcer
	NewYAMLFileSourcer          = config.NewYAMLFileSourcer
	ParseTOML                   = config.ParseTOML
	ParseYAML                   = config.ParseYAML
)
