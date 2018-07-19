package config

import (
	"os"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle/config/tag"
	. "github.com/onsi/gomega"
)

type (
	EnvConfigSuite struct{}
	TestConfigKey  struct{}
)

func (s *EnvConfigSuite) SetUpTest(t sweet.T) {
	os.Clearenv()
}

func (s *EnvConfigSuite) TestSimpleConfig(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
	Expect(chunk.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *EnvConfigSuite) TestNestedJSONDeserialization(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestEmbeddedJSONConfig{}
	)

	os.Setenv("APP_P1", `{"v_int": 3, "v_float": 3.14, "v_bool": true}`)
	os.Setenv("APP_P2", `{"v_int": 5, "v_float": 6.28, "v_bool": false}`)

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.P1).To(Equal(&TestJSONPayload{V1: 3, V2: 3.14, V3: true}))
	Expect(chunk.P2).To(Equal(&TestJSONPayload{V1: 5, V2: 6.28, V3: false}))
}

func (s *EnvConfigSuite) TestRequired(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestRequiredConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"no value supplied for field 'X'" +
		")",
	))
}

func (s *EnvConfigSuite) TestRequiredBadTag(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadRequiredConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"field 'X' has an invalid required tag" +
		")",
	))
}

func (s *EnvConfigSuite) TestDefault(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestDefaultConfig{}
	)

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *EnvConfigSuite) TestBadType(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("APP_X", "123") // silently converted to string
	os.Setenv("APP_Y", "foo")
	os.Setenv("APP_W", `bar`)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"value supplied for field 'Y' cannot be coerced into the expected type" +
		", " +
		"value supplied for field 'Z' cannot be coerced into the expected type" +
		")",
	))
}

func (s *EnvConfigSuite) TestBadDefaultType(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadDefaultConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"default value for field 'X' cannot be coerced into the expected type" +
		")",
	))
}

func (s *EnvConfigSuite) TestUnprefixed(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	os.Setenv("APP_Y", "456")

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(456))
}

func (s *EnvConfigSuite) TestNormalizedPrefix(t sweet.T) {
	var (
		config = NewEnvConfig("$foo-^-bar@")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("FOO_BAR_X", "foo")
	os.Setenv("FOO_BAR_Y", "123")

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
}

func (s *EnvConfigSuite) TestPostLoadConfig(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestPostLoadConfig{}
	)

	os.Setenv("APP_X", "3")
	Expect(config.Load(chunk)).To(BeNil())

	os.Setenv("APP_X", "-4")
	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"X must be positive" +
		")",
	))
}

func (s *EnvConfigSuite) TestUnsettableFields(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestUnsettableConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"field 'x' can not be set" +
		")",
	))
}

func (s *EnvConfigSuite) TestLoad(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
	Expect(chunk.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *EnvConfigSuite) TestLoadIsomorphicType(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
	Expect(chunk.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *EnvConfigSuite) TestLoadPostLoadWithConversion(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestPostLoadConversion{}
	)

	os.Setenv("APP_DURATION", "3")
	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.duration).To(Equal(time.Second * 3))
}

func (s *EnvConfigSuite) TestLoadPostLoadWithTags(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestPostLoadConversion{}
	)

	os.Setenv("APP_FOO_DURATION", "3")
	Expect(config.Load(chunk, tag.NewEnvTagPrefixer("foo"))).To(BeNil())
	Expect(chunk.duration).To(Equal(time.Second * 3))
}

func (s *EnvConfigSuite) TestBadConfigObjectTypes(t sweet.T) {
	Expect(NewEnvConfig("app").Load(nil)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"invalid type for configuration struct" +
		")",
	))

	Expect(NewEnvConfig("app").Load("foo")).To(MatchError("" +
		"failed to load config" +
		" (" +
		"invalid type for configuration struct" +
		")",
	))
}
