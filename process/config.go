package process

import (
	"errors"
	"time"
)

type (
	HTTPConfig struct {
		HTTPPort           int    `env:"http_port" default:"5000"`
		HTTPCertFile       string `env:"http_cert_file"`
		HTTPKeyFile        string `env:"http_key_file"`
		RawShutdownTimeout int    `env:"http_shutdown_timeout" default:"5"`

		ShutdownTimeout time.Duration
	}

	GRPCConfig struct {
		GRPCPort int `env:"grpc_port" default:"6000"`
	}

	WorkerConfig struct {
		RawWorkerTickInterval int `env:"worker_tick_interval" default:"0"`

		WorkerTickInterval time.Duration
	}

	httpConfigToken   struct{}
	grpcConfigToken   struct{}
	workerConfigToken struct{}
)

var (
	HTTPConfigToken   = httpConfigToken{}
	GRPCConfigToken   = grpcConfigToken{}
	WorkerConfigToken = workerConfigToken{}
	ErrBadCertConfig  = errors.New("cert file and key file must both be supplied or both be omitted")
)

func (c *HTTPConfig) PostLoad() error {
	if (c.HTTPCertFile == "") != (c.HTTPKeyFile == "") {
		return ErrBadCertConfig
	}

	c.ShutdownTimeout = time.Duration(c.RawShutdownTimeout) * time.Second
	return nil
}

func (c *WorkerConfig) PostLoad() error {
	c.WorkerTickInterval = time.Duration(c.RawWorkerTickInterval) * time.Second
	return nil
}
