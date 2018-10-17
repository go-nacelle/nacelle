package worker

import (
	"time"
)

type Config struct {
	RawWorkerTickInterval int `env:"worker_tick_interval" file:"worker_tick_interval" mask:"true" default:"0"`

	WorkerTickInterval time.Duration
}

func (c *Config) PostLoad() error {
	c.WorkerTickInterval = time.Duration(c.RawWorkerTickInterval) * time.Second
	return nil
}
