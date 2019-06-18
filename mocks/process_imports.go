package mocks

import (
	process "github.com/go-nacelle/process/mocks"
)

type (
	MockFinalizer        = process.MockFinalizer
	MockInitializer      = process.MockInitializer
	MockProcess          = process.MockProcess
	MockProcessContainer = process.MockProcessContainer
	MockRunner           = process.MockRunner
)

var (
	NewMockFinalizer        = process.NewMockFinalizer
	NewMockInitializer      = process.NewMockInitializer
	NewMockProcess          = process.NewMockProcess
	NewMockProcessContainer = process.NewMockProcessContainer
	NewMockRunner           = process.NewMockRunner
)
