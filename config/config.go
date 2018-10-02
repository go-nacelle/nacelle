package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/efritz/nacelle/config/tag"
)

type (
	// Config is a structure that can populate the exported fields of a
	// struct based on the value of the field `env` tags.
	Config interface {
		// Load populates a configuration object. The given tag modifiers
		// are applied to the configuration object pre-load. If the target
		// value conforms to the PostLoadConfig interface, the PostLoad
		// function may be called multiple times.
		Load(interface{}, ...tag.Modifier) error

		// MustInject calls Injects and panics on error.
		MustLoad(interface{}, ...tag.Modifier)
	}

	// PostLoadConfig is a marker interface for configuration objects
	// which should do some post-processing after being loaded. This
	// can perform additional casting (e.g. ints to time.Duration) and
	// more sophisticated validation (e.g. enum or exclusive values).
	PostLoadConfig interface {
		PostLoad() error
	}

	config struct {
		sourcer Sourcer
	}

	reflectField struct {
		field     reflect.Value
		fieldType reflect.StructField
	}
)

const (
	envTag      = "env"
	defaultTag  = "default"
	requiredTag = "required"
)

// NewConfig creates a config loader with the given sourcer.
func NewConfig(sourcer Sourcer) Config {
	return &config{
		sourcer: sourcer,
	}
}

func (c *config) Load(target interface{}, modifiers ...tag.Modifier) error {
	config, err := tag.ApplyModifiers(target, modifiers...)
	if err != nil {
		return err
	}

	errors := c.load(config)

	if len(errors) == 0 {
		sourceFields, _ := getExportedFields(config)
		targetFields, _ := getExportedFields(target)

		for i := 0; i < len(sourceFields); i++ {
			targetFields[i].field.Set(sourceFields[i].field)
		}

		if plc, ok := target.(PostLoadConfig); ok {
			if err := plc.PostLoad(); err != nil {
				errors = append(errors, err)
			}
		}
	}

	return loadError(errors)
}

// MustLoad calls Load and panics on error.
func (c *config) MustLoad(target interface{}, modifiers ...tag.Modifier) {
	if err := c.Load(target, modifiers...); err != nil {
		panic(err.Error())
	}
}

func (c *config) load(target interface{}) []error {
	objValue, objType, err := getIndirect(target)
	if err != nil {
		return []error{err}
	}

	return c.loadStruct(objValue, objType)
}

func (c *config) loadStruct(objValue reflect.Value, objType reflect.Type) []error {
	if objType.Kind() != reflect.Struct {
		return []error{fmt.Errorf(
			"invalid embedded type in configuration struct",
		)}
	}

	errors := []error{}
	for i := 0; i < objType.NumField(); i++ {
		var (
			field            = objValue.Field(i)
			fieldType        = objType.Field(i)
			defaultTagValue  = fieldType.Tag.Get(defaultTag)
			requiredTagValue = fieldType.Tag.Get(requiredTag)
		)

		if fieldType.Anonymous {
			errors = append(errors, c.loadStruct(field, fieldType.Type)...)
			continue
		}

		tagValues := []string{}
		for _, tag := range c.sourcer.Tags() {
			tagValues = append(tagValues, fieldType.Tag.Get(tag))
		}

		err := c.loadEnvField(
			field,
			fieldType.Name,
			tagValues,
			defaultTagValue,
			requiredTagValue,
		)

		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (c *config) loadEnvField(
	fieldValue reflect.Value,
	name string,
	tagValues []string,
	defaultTag string,
	requiredTag string,
) error {
	val, skip, ok := c.sourcer.Get(tagValues)
	if skip {
		return nil
	}

	if !fieldValue.IsValid() {
		return fmt.Errorf("field '%s' is invalid", name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field '%s' can not be set", name)
	}

	if ok {
		if !toJSON([]byte(val), fieldValue.Addr().Interface()) {
			return fmt.Errorf("value supplied for field '%s' cannot be coerced into the expected type", name)
		}

		return nil
	}

	if requiredTag != "" {
		val, err := strconv.ParseBool(requiredTag)
		if err != nil {
			return fmt.Errorf("field '%s' has an invalid required tag", name)
		}

		if val {
			return fmt.Errorf("no value supplied for field '%s'", name)
		}
	}

	if defaultTag != "" {
		if !fieldValue.IsValid() {
			return fmt.Errorf("field '%s' is invalid", name)
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("field '%s' can not be set", name)
		}

		if !toJSON([]byte(defaultTag), fieldValue.Addr().Interface()) {
			return fmt.Errorf("default value for field '%s' cannot be coerced into the expected type", name)
		}

		return nil
	}

	return nil
}

//
// Helpers

func loadError(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	messages := []string{}
	for _, err := range errors {
		messages = append(messages, err.Error())
	}

	return fmt.Errorf("failed to load config (%s)", strings.Join(messages, ", "))
}
