package nacelle

import "github.com/efritz/nacelle/config"

type (
	Config      = config.Config
	TagModifier = config.TagModifier
)

var (
	NewEnvConfig          = config.NewEnvConfig
	NewEnvTagPrefixer     = config.NewEnvTagPrefixer
	NewDefaultTagSetter   = config.NewDefaultTagSetter
	ApplyTagModifiers     = config.ApplyTagModifiers
	MustApplyTagModifiers = config.MustApplyTagModifiers
)
