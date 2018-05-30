package process

import (
	"fmt"
	"time"

	"github.com/efritz/nacelle"
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

// RegisterWorkerConfigs adds the required configs for a worker to the given map. If any tag
// modifiers are supplied, they are run over each of the required configs (this may require
// some knowledge about package internals).
func RegisterWorkerConfigs(m map[interface{}]interface{}, modifiers ...nacelle.TagModifier) error {
	c, err := nacelle.ApplyTagModifiers(&WorkerConfig{}, modifiers...)
	if err != nil {
		return err
	}

	m[WorkerConfigToken] = c
	return nil
}
