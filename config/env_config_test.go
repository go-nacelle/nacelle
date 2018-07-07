package config

import (
	"fmt"
	"os"
	"time"

	"github.com/aphistic/sweet"
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

	Expect(config.Register("simple", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	loadedChunk, err := config.Get("simple")
	Expect(err).To(BeNil())

	Expect(loadedChunk).To(BeIdenticalTo(chunk))
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

	Expect(config.Register("json-embedded", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	loadedChunk, err := config.Get("json-embedded")
	Expect(err).To(BeNil())

	Expect(loadedChunk).To(BeIdenticalTo(chunk))
	Expect(chunk.P1).To(Equal(&TestJSONPayload{V1: 3, V2: 3.14, V3: true}))
	Expect(chunk.P2).To(Equal(&TestJSONPayload{V1: 5, V2: 6.28, V3: false}))
}

func (s *EnvConfigSuite) TestRequired(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestRequiredConfig{}
	)

	Expect(config.Register("required-config", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("no value supplied for field 'X'")))
}

func (s *EnvConfigSuite) TestRequiredBadTag(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadRequiredConfig{}
	)

	Expect(config.Register("required-config", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("field 'X' has an invalid required tag")))
}

func (s *EnvConfigSuite) TestDefault(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestDefaultConfig{}
	)

	Expect(config.Register("default-config", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	loadedChunk, err := config.Get("default-config")
	Expect(err).To(BeNil())

	Expect(loadedChunk).To(BeIdenticalTo(chunk))
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

	Expect(config.Register("simple", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(2))
	Expect(errors).To(ContainElement(MatchError("value supplied for field 'Y' cannot be coerced into the expected type")))
	Expect(errors).To(ContainElement(MatchError("value supplied for field 'Z' cannot be coerced into the expected type")))
}

func (s *EnvConfigSuite) TestBadDefaultType(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadDefaultConfig{}
	)

	Expect(config.Register("bad-default", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("default value for field 'X' cannot be coerced into the expected type")))
}

func (s *EnvConfigSuite) TestUnprefixed(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	os.Setenv("APP_Y", "456")

	Expect(config.Register("simple", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	loadedChunk, err := config.Get("simple")
	Expect(err).To(BeNil())

	Expect(loadedChunk).To(BeIdenticalTo(chunk))
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

	Expect(config.Register("simple", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	loadedChunk, err := config.Get("simple")
	Expect(err).To(BeNil())

	Expect(loadedChunk).To(BeIdenticalTo(chunk))
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
}

func (s *EnvConfigSuite) TestToMap(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk1 = &TestSimpleConfig{}
		chunk2 = &TestEmbeddedJSONConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)
	os.Setenv("APP_P1", `{"v_int": 3, "v_float": 3.14, "v_bool": true}`)
	os.Setenv("APP_P2", `{"v_int": 5, "v_float": 6.28, "v_bool": false}`)

	config.MustRegister("simple", chunk1)
	config.MustRegister("json-embedded", chunk2)
	config.Load()

	dump, err := config.ToMap()
	Expect(err).To(BeNil())

	Expect(dump).To(HaveLen(5))
	Expect(dump["x"]).To(Equal("foo"))
	Expect(dump["y"]).To(Equal("123"))
	Expect(dump["q"]).To(MatchJSON(`["bar", "baz", "bonk"]`))
	Expect(dump["p1"]).To(MatchJSON(`{"v_int": 3, "v_float": 3.14, "v_bool": true}`))
	Expect(dump["p2"]).To(MatchJSON(`{"v_int": 5, "v_float": 6.28, "v_bool": false}`))
}

func (s *EnvConfigSuite) TestMask(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestMaskConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)

	config.MustRegister("masked", chunk)
	config.Load()

	dump, err := config.ToMap()
	Expect(err).To(BeNil())

	Expect(dump).To(HaveLen(1))
	Expect(dump["x"]).To(Equal("foo"))
}

func (s *EnvConfigSuite) TestBadMaskTag(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadMaskTagConfig{}
	)

	config.MustRegister("bad-mask-tag", chunk)
	config.Load()

	_, err := config.ToMap()
	Expect(err).To(MatchError("field 'X' has an invalid mask tag"))
}

func (s *EnvConfigSuite) TestPostLoadConfig(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestPostLoadConfig{}
	)

	Expect(config.Register("post-load", chunk)).To(BeNil())

	os.Setenv("APP_X", "3")
	Expect(config.Load()).To(BeEmpty())

	os.Setenv("APP_X", "-4")
	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("X must be positive")))
}

func (s *EnvConfigSuite) TestUnsettableFields(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestUnsettableConfig{}
	)

	Expect(config.Register("unsettable", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("field 'x' can not be set")))
}

func (s *EnvConfigSuite) TestRegistrationAfterLoad(t sweet.T) {
	config := NewEnvConfig("app")
	config.Register("pre-load", struct{}{})
	config.Load()
	Expect(config.Register("post-load", nil)).To(Equal(ErrAlreadyLoaded))
}

func (s *EnvConfigSuite) TestDuplicateRegistration(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	err1 := config.Register("dup", chunk)
	err2 := config.Register("dup", chunk)
	Expect(err1).To(BeNil())
	Expect(err2).To(MatchError("duplicate config key `dup`"))
}

func (s *EnvConfigSuite) TestMustRegisterPanics(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	config.MustRegister("dup", chunk)

	Expect(func() {
		config.MustRegister("dup", &TestSimpleConfig{})
	}).To(Panic())
}

func (s *EnvConfigSuite) TestGetBeforeLoad(t sweet.T) {
	_, err := NewEnvConfig("app").Get("pre-load")
	Expect(err).To(Equal(ErrNotLoaded))
}

func (s *EnvConfigSuite) TestGetUnregisteredKey(t sweet.T) {
	config := NewEnvConfig("app")
	config.Load()
	_, err := config.Get("unregistered")
	Expect(err).To(MatchError("unregistered config key `unregistered`"))
}

func (s *EnvConfigSuite) TestMustGetPanics(t sweet.T) {
	Expect(func() {
		NewEnvConfig("app").MustGet("unregistered")
	}).To(Panic())
}

func (s *EnvConfigSuite) TestFetch(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)

	Expect(config.Register("simple", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	target := &TestSimpleConfig{}
	Expect(config.Fetch("simple", target)).To(BeNil())
	Expect(target.X).To(Equal("foo"))
	Expect(target.Y).To(Equal(123))
	Expect(target.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *EnvConfigSuite) TestFetchIsomorphicType(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)

	Expect(config.Register("simple", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	target := &TestSimpleConfigClone{}
	Expect(config.Fetch("simple", target)).To(BeNil())
	Expect(target.X).To(Equal("foo"))
	Expect(target.Y).To(Equal(123))
	Expect(target.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *EnvConfigSuite) TestFetchBadType(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	os.Setenv("APP_X", "foo")
	os.Setenv("APP_Y", "123")
	os.Setenv("APP_W", `["bar", "baz", "bonk"]`)

	Expect(config.Register("simple", chunk)).To(BeNil())
	Expect(config.Load()).To(BeEmpty())

	target := &TestDefaultConfig{}
	Expect(config.Fetch("simple", target)).NotTo(BeNil())
}

func (s *EnvConfigSuite) TestFetchPostLoadWithConversion(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestPostLoadConversion{}
	)

	Expect(config.Register("post-load", chunk)).To(BeNil())

	os.Setenv("APP_DURATION", "3")
	Expect(config.Load()).To(BeEmpty())
	Expect(chunk.duration).To(Equal(time.Second * 3))

	target := &TestPostLoadConversion{}
	Expect(config.Fetch("post-load", target)).To(BeNil())
	Expect(target.duration).To(Equal(time.Second * 3))
}

func (s *EnvConfigSuite) TestFetchWithConfigTagRoundtrip(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = MustApplyTagModifiers(&TestPostLoadConversion{}, NewEnvTagPrefixer("foo"))
	)

	Expect(config.Register("post-load", chunk)).To(BeNil())

	os.Setenv("APP_FOO_DURATION", "3")
	Expect(config.Load()).To(BeEmpty())

	target := &TestPostLoadConversion{}
	Expect(config.Fetch("post-load", target)).To(BeNil())
	Expect(target.duration).To(Equal(time.Second * 3))
}

func (s *EnvConfigSuite) TestSerializeKey(t sweet.T) {
	Expect(serializeKey("foo")).To(Equal("foo"))
	Expect(serializeKey(TestStringer{})).To(Equal("bar"))
	Expect(serializeKey(TestConfigKey{})).To(Equal("TestConfigKey"))
	Expect(serializeKey(&TestConfigKey{})).To(Equal("TestConfigKey"))
}

//
// Chunks

type (
	TestSimpleConfig struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"q"`
	}

	TestSimpleConfigClone struct {
		X string
		Y int
		Z []string
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

//
// Helpers

type TestStringer struct{}

func (TestStringer) String() string {
	return "bar"
}
