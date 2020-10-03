package mocks

import config "github.com/go-nacelle/config/mocks"

type (
	MockConfig      = config.MockConfig
	MockSourcer     = config.MockSourcer
	MockTagModifier = config.MockTagModifier
)

var (
	NewMockConfig      = config.NewMockConfig
	NewMockSourcer     = config.NewMockSourcer
	NewMockTagModifier = config.NewMockTagModifier
)
