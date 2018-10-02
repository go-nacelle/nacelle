package nacelle

import (
	"github.com/efritz/nacelle/config"
	"github.com/efritz/nacelle/config/tag"
)

type (
	Config        = config.Config
	ConfigSourcer = config.Sourcer
	TagModifier   = tag.Modifier
)

var (
	NewConfig              = config.NewConfig
	NewEnvSourcer          = config.NewEnvSourcer
	NewFileSourcer         = config.NewFileSourcer
	NewOptionalFileSourcer = config.NewOptionalFileSourcer
	NewYAMLFileSourcer     = config.NewYAMLFileSourcer
	NewTOMLFileSourcer     = config.NewTOMLFileSourcer
	NewMultiSourcer        = config.NewMultiSourcer
	NewEnvTagPrefixer      = tag.NewEnvTagPrefixer
	NewDefaultTagSetter    = tag.NewDefaultTagSetter
)
