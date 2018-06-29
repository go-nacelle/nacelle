package process

import (
	"sync"
	"time"
)

type (
	initializerMeta struct {
		Initializer
		name    string
		timeout time.Duration
	}

	processMeta struct {
		Process
		name        string
		priority    int
		silentExit  bool
		initTimeout time.Duration
		once        *sync.Once
	}
)

func (m *initializerMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}

func (m *processMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}

func newInitializerMeta(initializer Initializer) *initializerMeta {
	return &initializerMeta{
		Initializer: initializer,
	}
}

func newProcessMeta(process Process) *processMeta {
	return &processMeta{
		Process: process,
		once:    &sync.Once{},
	}
}
