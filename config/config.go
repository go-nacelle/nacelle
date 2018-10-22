package config

import "github.com/efritz/zubrin"

type (
	Config        = zubrin.Config
	LoggingConfig = zubrin.LoggingConfig
	Sourcer       = zubrin.Sourcer
	TagModifier   = zubrin.TagModifier
)

var (
	NewConfig                   = zubrin.NewConfig
	NewEnvSourcer               = zubrin.NewEnvSourcer
	NewFileSourcer              = zubrin.NewFileSourcer
	NewGlobSourcer              = zubrin.NewGlobSourcer
	NewOptionalFileSourcer      = zubrin.NewOptionalFileSourcer
	NewDirectorySourcer         = zubrin.NewDirectorySourcer
	NewOptionalDirectorySourcer = zubrin.NewOptionalDirectorySourcer
	NewYAMLFileSourcer          = zubrin.NewYAMLFileSourcer
	NewTOMLFileSourcer          = zubrin.NewTOMLFileSourcer
	NewMultiSourcer             = zubrin.NewMultiSourcer
	NewEnvTagPrefixer           = zubrin.NewEnvTagPrefixer
	NewFileTagPrefixer          = zubrin.NewFileTagPrefixer
	NewDefaultTagSetter         = zubrin.NewDefaultTagSetter
)
