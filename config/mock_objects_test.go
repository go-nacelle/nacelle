package config

import (
	"fmt"
	"time"
)

type (
	TestSimpleConfig struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	TestEmbeddedJSONConfig struct {
		P1 *TestJSONPayload `env:"p1"`
		P2 *TestJSONPayload `env:"p2"`
	}

	TestJSONPayload struct {
		V1 int     `json:"v_int"`
		V2 float64 `json:"v_float"`
		V3 bool    `json:"v_bool"`
	}

	TestRequiredConfig struct {
		X string `env:"x" required:"true"`
	}

	TestBadRequiredConfig struct {
		X string `env:"x" required:"yup"`
	}

	TestDefaultConfig struct {
		X string   `env:"x" default:"foo"`
		Y []string `env:"y" default:"[\"bar\", \"baz\", \"bonk\"]"`
	}

	TestBadDefaultConfig struct {
		X int `env:"x" default:"foo"`
	}

	TestUnsettableConfig struct {
		x int `env:"s"`
	}

	TestPostLoadConfig struct {
		X int `env:"X"`
	}

	TestPostLoadConversion struct {
		RawDuration int `env:"duration"`
		duration    time.Duration
	}

	TestMaskConfig struct {
		X string   `env:"x"`
		Y int      `env:"y" mask:"true"`
		Z []string `env:"w" mask:"true"`
	}

	TestBadMaskTagConfig struct {
		X string `env:"x" mask:"34"`
	}

	TestParentConfig struct {
		ChildConfig
		X int `env:"x"`
		Y int `env:"y"`
	}

	ChildConfig struct {
		A int `env:"a"`
		B int `env:"b"`
		C int `env:"c"`
	}

	TestBadParentConfig struct {
		*ChildConfig
		X int `env:"x"`
		Y int `env:"y"`
	}
)

func (c *TestPostLoadConfig) PostLoad() error {
	if c.X < 0 {
		return fmt.Errorf("X must be positive")
	}

	return nil
}

func (c *TestPostLoadConversion) PostLoad() error {
	c.duration = time.Duration(c.RawDuration) * time.Second
	return nil
}

func (c *ChildConfig) PostLoad() error {
	if c.A >= c.B || c.B >= c.C {
		return fmt.Errorf("fields must be increasing")
	}

	return nil
}

//
// Helpers

type TestStringer struct{}

func (TestStringer) String() string {
	return "bar"
}
