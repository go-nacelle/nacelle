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
	Config interface {
		Load() []error
		Register(key interface{}, config interface{}) error
		MustRegister(key interface{}, config interface{})
		Get(key interface{}) (interface{}, error)
		MustGet(key interface{}) interface{}
		ToMap() (map[string]interface{}, error)
	}

	PostLoadConfig interface {
		PostLoad() error
	}

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
)

var (
	ErrAlreadyLoaded         = errors.New("config already loaded")
	ErrNotLoaded             = errors.New("config not loaded")
	ErrDuplicateConfigKey    = errors.New("duplicate config key")
	ErrUnregisteredConfigKey = errors.New("no config registered to key")

	replacer = strings.NewReplacer(
		"\n", `\n`,
		"\t", `\t`,
		"\r", `\r`,
	)
)

func NewEnvConfig(prefix string) Config {
	return &EnvConfig{
		prefix: prefix,
		chunks: map[interface{}]interface{}{},
	}
}

func (c *EnvConfig) Register(key interface{}, config interface{}) error {
	if c.loaded {
		return ErrAlreadyLoaded
	}

	if _, ok := c.chunks[key]; ok {
		return ErrDuplicateConfigKey
	}

	c.chunks[key] = config
	return nil
}

func (c *EnvConfig) MustRegister(key interface{}, config interface{}) {
	if err := c.Register(key, config); err != nil {
		panic(err.Error())
	}
}

func (c *EnvConfig) Get(key interface{}) (interface{}, error) {
	if !c.loaded {
		return nil, ErrNotLoaded
	}

	if config, ok := c.chunks[key]; ok {
		return config, nil
	}

	return nil, ErrUnregisteredConfigKey
}

func (c *EnvConfig) MustGet(key interface{}) interface{} {
	config, err := c.Get(key)
	if err != nil {
		panic(err.Error())
	}

	return config
}

func (c *EnvConfig) Load() []error {
	c.loaded = true
	errors := []error{}

	for _, chunk := range c.chunks {
		errors = loadChunk(chunk, errors, c.prefix)
	}

	return errors
}

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
			fieldType   = ot.Field(i)
			fieldValue  = oi.Field(i)
			envTag      = fieldType.Tag.Get(envTag)
			defaultTag  = fieldType.Tag.Get(defaultTag)
			requiredTag = fieldType.Tag.Get(requiredTag)
		)

		if envTag == "" {
			continue
		}

		envTags := []string{
			strings.ToUpper(fmt.Sprintf("%s_%s", prefix, envTag)),
			strings.ToUpper(envTag),
		}

		err := loadEnvField(
			fieldType,
			fieldValue,
			envTags,
			defaultTag,
			requiredTag,
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
			fieldType  = ot.Field(i)
			fieldValue = oi.Field(i)
			envTag     = fieldType.Tag.Get(envTag)
			maskTag    = fieldType.Tag.Get(maskTag)
		)

		if envTag == "" {
			continue
		}

		if maskTag != "" {
			val, err := strconv.ParseBool(maskTag)
			if err != nil {
				return fmt.Errorf("field '%s' has an invalid mask tag", fieldType.Name)
			}

			if val {
				continue
			}
		}

		if fieldValue.Kind() == reflect.String {
			m[fieldType.Name] = fmt.Sprintf("%s", fieldValue)
		} else {
			data, err := json.Marshal(fieldValue.Interface())
			if err != nil {
				return err
			}

			m[fieldType.Name] = string(data)
		}
	}

	return nil
}
