package process

import (
	"fmt"
	"time"
)

type (
	WorkerConfig struct {
		RawWorkerTickInterval int `env:"worker_tick_interval" default:"0"`

		WorkerTickInterval time.Duration
	}

	workerConfigToken string
)

var WorkerConfigToken = MakeWorkerConfigToken("default")

func MakeWorkerConfigToken(name string) interface{} {
	return workerConfigToken(fmt.Sprintf("nacelle-process-worker-%s", name))
}

func (c *WorkerConfig) PostLoad() error {
	c.WorkerTickInterval = time.Duration(c.RawWorkerTickInterval) * time.Second
	return nil
}
