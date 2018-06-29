package worker

import (
	"fmt"
	"time"
)

type (
	Config struct {
		RawWorkerTickInterval int `env:"worker_tick_interval" default:"0"`

		WorkerTickInterval time.Duration
	}

	configToken string
)

var ConfigToken = NewConfigToken("default")

func NewConfigToken(name string) interface{} {
	return configToken(fmt.Sprintf("nacelle-base-worker-%s", name))
}

func (c *Config) PostLoad() error {
	c.WorkerTickInterval = time.Duration(c.RawWorkerTickInterval) * time.Second
	return nil
}
