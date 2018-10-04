package config

import "github.com/efritz/zubrin"

type (
	Config        = zubrin.Config
	LoggingConfig = zubrin.LoggingConfig
	Sourcer       = zubrin.Sourcer
	TagModifier   = zubrin.TagModifier
)

var (
	NewConfig              = zubrin.NewConfig
	NewEnvSourcer          = zubrin.NewEnvSourcer
	NewFileSourcer         = zubrin.NewFileSourcer
	NewOptionalFileSourcer = zubrin.NewOptionalFileSourcer
	NewYAMLFileSourcer     = zubrin.NewYAMLFileSourcer
	NewTOMLFileSourcer     = zubrin.NewTOMLFileSourcer
	NewMultiSourcer        = zubrin.NewMultiSourcer
	NewEnvTagPrefixer      = zubrin.NewEnvTagPrefixer
	NewDefaultTagSetter    = zubrin.NewDefaultTagSetter
)
