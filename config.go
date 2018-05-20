package nacelle

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type (
	// Config is a structure that maintains chunks of an application's
	// configuration values accessible via arbitrary key. This design is
	// meant for a library-driven architecture so that multiple pieces
	// of an application can register their own configuration requirements
	// independently from the core.
	Config interface {
		// Load attempts to read an external environment for values
		// and modifies the config's internal state on success.
		Load() []error

		// Register associates an empty configuration object. This key should
		// be unique to the application as it is generally error to register
		// the same key twice.
		Register(key interface{}, config interface{}) error

		// MustRegister calls Register and panics on error.
		MustRegister(key interface{}, config interface{})

		// Get retrieves a configuration object by its key. It is an error
		// to request a non-registered key or to Get before a call to Load.
		Get(key interface{}) (interface{}, error)

		// MustGet calls Get and panics on error.
		MustGet(key interface{}) interface{}

		// ToMap will convert the configuration values into a printable
		// or loggable map.
		ToMap() (map[string]interface{}, error)
	}

	// PostLoadConfig is a marker interface for configuration objects
	// which should do some post-processing after being loaded. This
	// can perform additional casting (e.g. ints to time.Duration) and
	// more sophisticated validation (e.g. enum or exclusive values).
	PostLoadConfig interface {
		PostLoad() error
	}

	// EnvConfig is a Config object that reads from the OS environment.
	EnvConfig struct {
		prefix string
		chunks map[interface{}]interface{}
		loaded bool
	}
)

const (
	envTag      = "env"
	maskTag     = "mask"
	defaultTag  = "default"
	requiredTag = "required"
	displayTag  = "display"
)

var (
	// ErrAlreadyLoaded is returned on a second call to Config#Load.
	ErrAlreadyLoaded = errors.New("config already loaded")

	// ErrNotLoaded is returned on a call to Get without first calling Load.
	ErrNotLoaded = errors.New("config not loaded")

	replacer = strings.NewReplacer(
		"\n", `\n`,
		"\t", `\t`,
		"\r", `\r`,
	)
)

// NewEnvConfig creates a EnvConfig object with the given prefix. If supplied,
// the {PREFIX}{NAME} envvar is read before falling back to the {NAME} envvar.
func NewEnvConfig(prefix string) Config {
	return &EnvConfig{
		prefix: prefix,
		chunks: map[interface{}]interface{}{},
	}
}

// Register associates a zero-valued struct whose exported fields should be tagged
// as `env:"name"` with a key. It is an error to register the same key twice.
func (c *EnvConfig) Register(key interface{}, config interface{}) error {
	if c.loaded {
		return ErrAlreadyLoaded
	}

	if _, ok := c.chunks[key]; ok {
		return fmt.Errorf("duplicate config key `%s`", serializeKey(key))
	}

	c.chunks[key] = config
	return nil
}

// MustRegister calls Register and panics on error.
func (c *EnvConfig) MustRegister(key interface{}, config interface{}) {
	if err := c.Register(key, config); err != nil {
		panic(err.Error())
	}
}

// Get retrieves the populated struct by its key.
func (c *EnvConfig) Get(key interface{}) (interface{}, error) {
	if !c.loaded {
		return nil, ErrNotLoaded
	}

	if config, ok := c.chunks[key]; ok {
		return config, nil
	}

	return nil, fmt.Errorf("unregistered config key `%s`", serializeKey(key))
}

// MustGet calls Get and panics on error.
func (c *EnvConfig) MustGet(key interface{}) interface{} {
	config, err := c.Get(key)
	if err != nil {
		panic(err.Error())
	}

	return config
}

// Load each registered struct with values from the environment. If a struct field
// is tagged as `required:"true"` and no value (nor default value) is supplied, an
// error is generated. If a struct field is tagged with a `default:"value"` value and
// no value is supplied from the environment, that value is used as if it came from
// the environment. The values that are pulled from the environment are attempted to
// be treated as JSON and, on failure, are treated as a string before assigning them
// to registered struct fields. This allows lists and map types to be expressed easily.
func (c *EnvConfig) Load() []error {
	c.loaded = true

	errors := []error{}
	for _, chunk := range c.chunks {
		errors = loadChunk(chunk, errors, c.prefix)
	}

	return errors
}

// ToMap will serialize the loaded config structs into a map. If a struct field has a
// `mask:"true"` tag it will be omitted form the result. If a struct field has the tag
// `display:"name"`, then the tag's value will be used in place of the field name.
func (c *EnvConfig) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}

	for _, chunk := range c.chunks {
		if err := dumpChunk(chunk, m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func loadChunk(obj interface{}, errors []error, prefix string) []error {
	var (
		ov = reflect.ValueOf(obj)
		oi = reflect.Indirect(ov)
		ot = oi.Type()
	)

	for i := 0; i < ot.NumField(); i++ {
		var (
			fieldType        = ot.Field(i)
			fieldValue       = oi.Field(i)
			envTagValue      = fieldType.Tag.Get(envTag)
			defaultTagValue  = fieldType.Tag.Get(defaultTag)
			requiredTagValue = fieldType.Tag.Get(requiredTag)
		)

		if envTagValue == "" {
			continue
		}

		envTags := []string{
			strings.ToUpper(fmt.Sprintf("%s_%s", prefix, envTagValue)),
			strings.ToUpper(envTagValue),
		}

		err := loadEnvField(
			fieldType,
			fieldValue,
			envTags,
			defaultTagValue,
			requiredTagValue,
		)

		if err != nil {
			errors = append(errors, err)
		}
	}

	if plc, ok := obj.(PostLoadConfig); ok {
		if err := plc.PostLoad(); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func loadEnvField(fieldType reflect.StructField, fieldValue reflect.Value, envTags []string, defaultTag, requiredTag string) error {
	if !fieldValue.IsValid() {
		return fmt.Errorf("field '%s' is invalid", fieldType.Name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field '%s' can not be set", fieldType.Name)
	}

	val, ok := getFirst(envTags)
	if ok {
		if !toJSON([]byte(val), fieldValue.Addr().Interface()) {
			return fmt.Errorf("value supplied for field '%s' cannot be coerced into the expected type", fieldType.Name)
		}

		return nil
	}

	if requiredTag != "" {
		val, err := strconv.ParseBool(requiredTag)
		if err != nil {
			return fmt.Errorf("field '%s' has an invalid required tag", fieldType.Name)
		}

		if val {
			return fmt.Errorf("no value supplied for field '%s'", fieldType.Name)
		}
	}

	if defaultTag != "" {
		if !toJSON([]byte(defaultTag), fieldValue.Addr().Interface()) {
			return fmt.Errorf("default value for field '%s' cannot be coerced into the expected type", fieldType.Name)
		}

		return nil
	}

	return nil
}

func getFirst(envTags []string) (string, bool) {
	for _, envTag := range envTags {
		if val, ok := os.LookupEnv(envTag); ok {
			return val, ok
		}
	}

	return "", false
}

func toJSON(data []byte, v interface{}) bool {
	if json.Unmarshal(data, v) == nil {
		return true
	}

	ptr := reflect.ValueOf(v)

	if ptr.Kind() == reflect.Ptr && reflect.Indirect(ptr).Kind() == reflect.String {
		if json.Unmarshal(quoteJSON(data), v) == nil {
			return true
		}
	}

	return false
}

func quoteJSON(data []byte) []byte {
	return []byte(fmt.Sprintf(`"%s"`, replacer.Replace(string(data))))
}

func dumpChunk(obj interface{}, m map[string]interface{}) error {
	var (
		ov = reflect.ValueOf(obj)
		oi = reflect.Indirect(ov)
		ot = oi.Type()
	)

	for i := 0; i < ot.NumField(); i++ {
		var (
			fieldType       = ot.Field(i)
			fieldValue      = oi.Field(i)
			envTagValue     = fieldType.Tag.Get(envTag)
			maskTagValue    = fieldType.Tag.Get(maskTag)
			displayTagValue = fieldType.Tag.Get(displayTag)
			displayName     = ""
		)

		if displayTagValue != "" {
			displayName = displayTagValue
		} else {
			if envTagValue == "" {
				continue
			}

			displayName = strings.ToLower(envTagValue)
		}

		if maskTagValue != "" {
			val, err := strconv.ParseBool(maskTagValue)
			if err != nil {
				return fmt.Errorf("field '%s' has an invalid mask tag", fieldType.Name)
			}

			if val {
				continue
			}
		}

		if fieldValue.Kind() == reflect.String {
			m[displayName] = fmt.Sprintf("%s", fieldValue)
		} else {
			data, err := json.Marshal(fieldValue.Interface())
			if err != nil {
				return err
			}

			m[displayName] = string(data)
		}
	}

	return nil
}
