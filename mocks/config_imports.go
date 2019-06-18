package mocks

import (
	config "github.com/go-nacelle/config/mocks"
)

type (
	MockConfig      = config.MockConfig
	MockLogger      = config.MockLogger
	MockSourcer     = config.MockSourcer
	MockTagModifier = config.MockTagModifier
)

var (
	NewMockConfig      = config.NewMockConfig
	NewMockLogger      = config.NewMockLogger
	NewMockSourcer     = config.NewMockSourcer
	NewMockTagModifier = config.NewMockTagModifier
)
