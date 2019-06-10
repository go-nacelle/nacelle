package mocks

import (
	config "github.com/go-nacelle/config/mocks"
	process "github.com/go-nacelle/process/mocks"
	service "github.com/go-nacelle/service/mocks"
)

type (
	MockFinalizer        = process.MockFinalizer
	MockConfig           = config.MockConfig
	MockInitializer      = process.MockInitializer
	MockLogger           = config.MockLogger
	MockProcess          = process.MockProcess
	MockProcessContainer = process.MockProcessContainer
	MockRunner           = process.MockRunner
	MockServiceContainer = service.MockServiceContainer
	MockSourcer          = config.MockSourcer
	MockTagModifier      = config.MockTagModifier
)

var (
	NewMockFinalizer        = process.NewMockFinalizer
	NewMockInitializer      = process.NewMockInitializer
	NewMockProcess          = process.NewMockProcess
	NewMockProcessContainer = process.NewMockProcessContainer
	NewMockRunner           = process.NewMockRunner
	NewMockConfig           = config.NewMockConfig
	NewMockLogger           = config.NewMockLogger
	NewMockServiceContainer = service.NewMockServiceContainer
	NewMockSourcer          = config.NewMockSourcer
	NewMockTagModifier      = config.NewMockTagModifier
)
