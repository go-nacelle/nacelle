package nacelle

import (
	"github.com/efritz/nacelle/config"
	"github.com/efritz/nacelle/config/tag"
)

type (
	Config      = config.Config
	TagModifier = tag.Modifier
)

var (
	NewEnvConfig        = config.NewEnvConfig
	NewEnvTagPrefixer   = tag.NewEnvTagPrefixer
	NewDefaultTagSetter = tag.NewDefaultTagSetter
)
