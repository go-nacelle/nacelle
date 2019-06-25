package mocks

import (
	log "github.com/go-nacelle/log/mocks"
)

type (
	MockLogger = log.MockLogger
)

var (
	NewMockLogger = log.NewMockLogger
)
