package config

import (
	"os"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type EnvSourcerSuite struct{}

func (s *EnvSourcerSuite) TestUnprefixed(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	os.Setenv("APP_Y", "456")

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(456))
}

func (s *EnvSourcerSuite) TestNormalizedPrefix(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("$foo-^-bar@"))
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("FOO_BAR_X", "foo")
	os.Setenv("FOO_BAR_Y", "123")

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
}
