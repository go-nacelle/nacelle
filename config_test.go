package nacelle

import (
	"errors"
	"os"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) SetUpTest(t sweet.T) {
	os.Clearenv()
}

func (s *ConfigSuite) TestSimpleConfig(t sweet.T) {
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

func (s *ConfigSuite) TestNestedJSONDeserialization(t sweet.T) {
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

func (s *ConfigSuite) TestRequired(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestRequiredConfig{}
	)

	Expect(config.Register("required-config", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("no value supplied for field 'X'")))
}

func (s *ConfigSuite) TestRequiredBadTag(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadRequiredConfig{}
	)

	Expect(config.Register("required-config", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("field 'X' has an invalid required tag")))
}

func (s *ConfigSuite) TestDefault(t sweet.T) {
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

func (s *ConfigSuite) TestBadType(t sweet.T) {
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

func (s *ConfigSuite) TestBadDefaultType(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadDefaultConfig{}
	)

	Expect(config.Register("bad-default", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("default value for field 'X' cannot be coerced into the expected type")))
}

func (s *ConfigSuite) TestUnprefixed(t sweet.T) {
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

func (s *ConfigSuite) TestToMap(t sweet.T) {
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
	Expect(dump["X"]).To(Equal("foo"))
	Expect(dump["Y"]).To(Equal("123"))
	Expect(dump["Z"]).To(MatchJSON(`["bar", "baz", "bonk"]`))
	Expect(dump["P1"]).To(MatchJSON(`{"v_int": 3, "v_float": 3.14, "v_bool": true}`))
	Expect(dump["P2"]).To(MatchJSON(`{"v_int": 5, "v_float": 6.28, "v_bool": false}`))
}

func (s *ConfigSuite) TestMask(t sweet.T) {
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
	Expect(dump["X"]).To(Equal("foo"))
}

func (s *ConfigSuite) TestBadMaskTag(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestBadMaskTagConfig{}
	)

	config.MustRegister("bad-mask-tag", chunk)
	config.Load()

	_, err := config.ToMap()
	Expect(err).To(MatchError("field 'X' has an invalid mask tag"))
}

func (s *ConfigSuite) TestPostLoadConfig(t sweet.T) {
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

func (s *ConfigSuite) TestUnsettableFields(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestUnsettableConfig{}
	)

	Expect(config.Register("unsettable", chunk)).To(BeNil())

	errors := config.Load()
	Expect(errors).To(HaveLen(1))
	Expect(errors).To(ContainElement(MatchError("field 'x' can not be set")))
}

func (s *ConfigSuite) TestRegistrationAfterLoad(t sweet.T) {
	config := NewEnvConfig("app")
	config.Register("pre-load", struct{}{})
	config.Load()
	Expect(config.Register("post-load", nil)).To(Equal(ErrAlreadyLoaded))
}

func (s *ConfigSuite) TestDuplicateRegistration(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	err1 := config.Register("dup", chunk)
	err2 := config.Register("dup", chunk)
	Expect(err1).To(BeNil())
	Expect(err2).To(Equal(ErrDuplicateConfigKey))
}

func (s *ConfigSuite) TestMustRegisterPanics(t sweet.T) {
	var (
		config = NewEnvConfig("app")
		chunk  = &TestSimpleConfig{}
	)

	config.MustRegister("dup", chunk)

	Expect(func() {
		config.MustRegister("dup", &TestSimpleConfig{})
	}).To(Panic())
}

func (s *ConfigSuite) TestGetBeforeLoad(t sweet.T) {
	_, err := NewEnvConfig("app").Get("pre-load")
	Expect(err).To(Equal(ErrNotLoaded))
}

func (s *ConfigSuite) TestGetUnregisteredKey(t sweet.T) {
	config := NewEnvConfig("app")
	config.Load()
	_, err := config.Get("unregistered")
	Expect(err).To(Equal(ErrUnregisteredConfigKey))
}

func (s *ConfigSuite) TestMustGetPanics(t sweet.T) {
	Expect(func() {
		NewEnvConfig("app").MustGet("unregistered")
	}).To(Panic())
}

//
// Chunks

type (
	TestSimpleConfig struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w"`
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
		return errors.New("X must be positive")
	}

	return nil
}
