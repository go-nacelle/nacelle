package config

import (
	"os"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/nacelle/config/tag"
)

type ConfigSuite struct{}

func (s *ConfigSuite) SetUpTest(t sweet.T) {
	os.Clearenv()
}

func (s *ConfigSuite) TestSimpleConfig(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
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

func (s *ConfigSuite) TestNestedJSONDeserialization(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestEmbeddedJSONConfig{}
	)

	os.Setenv("APP_P1", `{"v_int": 3, "v_float": 3.14, "v_bool": true}`)
	os.Setenv("APP_P2", `{"v_int": 5, "v_float": 6.28, "v_bool": false}`)

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.P1).To(Equal(&TestJSONPayload{V1: 3, V2: 3.14, V3: true}))
	Expect(chunk.P2).To(Equal(&TestJSONPayload{V1: 5, V2: 6.28, V3: false}))
}

func (s *ConfigSuite) TestRequired(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestRequiredConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"no value supplied for field 'X'" +
		")",
	))
}

func (s *ConfigSuite) TestRequiredBadTag(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestBadRequiredConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"field 'X' has an invalid required tag" +
		")",
	))
}

func (s *ConfigSuite) TestDefault(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestDefaultConfig{}
	)

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *ConfigSuite) TestBadType(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
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

func (s *ConfigSuite) TestBadDefaultType(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestBadDefaultConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"default value for field 'X' cannot be coerced into the expected type" +
		")",
	))
}

func (s *ConfigSuite) TestPostLoadConfig(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
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

func (s *ConfigSuite) TestUnsettableFields(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestUnsettableConfig{}
	)

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"field 'x' can not be set" +
		")",
	))
}

func (s *ConfigSuite) TestLoad(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
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

func (s *ConfigSuite) TestLoadIsomorphicType(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
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

func (s *ConfigSuite) TestLoadPostLoadWithConversion(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestPostLoadConversion{}
	)

	os.Setenv("APP_DURATION", "3")
	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.duration).To(Equal(time.Second * 3))
}

func (s *ConfigSuite) TestLoadPostLoadWithTags(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestPostLoadConversion{}
	)

	os.Setenv("APP_FOO_DURATION", "3")
	Expect(config.Load(chunk, tag.NewEnvTagPrefixer("foo"))).To(BeNil())
	Expect(chunk.duration).To(Equal(time.Second * 3))
}

func (s *ConfigSuite) TestBadConfigObjectTypes(t sweet.T) {
	Expect(NewConfig(NewEnvSourcer("app")).Load(nil)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"configuration target is not a pointer to struct" +
		")",
	))

	Expect(NewConfig(NewEnvSourcer("app")).Load("foo")).To(MatchError("" +
		"failed to load config" +
		" (" +
		"configuration target is not a pointer to struct" +
		")",
	))
}

func (s *ConfigSuite) TestEmbeddedConfig(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestParentConfig{}
	)

	os.Setenv("APP_A", "1")
	os.Setenv("APP_B", "2")
	os.Setenv("APP_C", "3")
	os.Setenv("APP_X", "4")
	os.Setenv("APP_Y", "5")

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal(4))
	Expect(chunk.Y).To(Equal(5))
	Expect(chunk.A).To(Equal(1))
	Expect(chunk.B).To(Equal(2))
	Expect(chunk.C).To(Equal(3))
}

func (s *ConfigSuite) TestEmbeddedConfigWithTags(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestParentConfig{}
	)

	os.Setenv("APP_FOO_A", "1")
	os.Setenv("APP_FOO_B", "2")
	os.Setenv("APP_FOO_C", "3")
	os.Setenv("APP_FOO_X", "4")
	os.Setenv("APP_FOO_Y", "5")

	Expect(config.Load(chunk, tag.NewEnvTagPrefixer("foo"))).To(BeNil())
	Expect(chunk.X).To(Equal(4))
	Expect(chunk.Y).To(Equal(5))
	Expect(chunk.A).To(Equal(1))
	Expect(chunk.B).To(Equal(2))
	Expect(chunk.C).To(Equal(3))
}

func (s *ConfigSuite) TestEmbeddedConfigPostLoad(t sweet.T) {
	var (
		config = NewConfig(NewEnvSourcer("app"))
		chunk  = &TestParentConfig{}
	)

	os.Setenv("APP_A", "1")
	os.Setenv("APP_B", "3")
	os.Setenv("APP_C", "2")
	os.Setenv("APP_X", "4")
	os.Setenv("APP_Y", "5")

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"fields must be increasing" +
		")",
	))
}

func (s *ConfigSuite) TestBadEmbeddedObjectType(t sweet.T) {
	Expect(NewConfig(NewEnvSourcer("app")).Load(&TestBadParentConfig{})).To(MatchError("" +
		"failed to load config" +
		" (" +
		"invalid embedded type in configuration struct" +
		")",
	))
}
