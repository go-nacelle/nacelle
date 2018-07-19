package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/efritz/nacelle/config/tag"
)

type (
	// envConfig is a Config object that reads from the OS environment.
	envConfig struct {
		prefix string
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

var (
	replacer = strings.NewReplacer(
		"\n", `\n`,
		"\t", `\t`,
		"\r", `\r`,
	)

	replacePattern = regexp.MustCompile(`[^A-Za-z0-9_]+`)
)

// NewEnvConfig creates a EnvConfig object with the given prefix. If supplied,
// the {PREFIX}_{NAME} envvar is read before falling back to the {NAME} envvar.
// The prefix will be normalized (replaces all non-alpha characters with an
// underscore and trims leading, trailing, and collapses consecutive underscores).
func NewEnvConfig(prefix string) Config {
	normalizedPrefix := strings.Trim(
		string(replacePattern.ReplaceAll(
			[]byte(prefix),
			[]byte("_"),
		)),
		"_",
	)

	return &envConfig{
		prefix: normalizedPrefix,
	}
}

func (c *envConfig) Load(target interface{}, modifiers ...tag.Modifier) error {
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
func (c *envConfig) MustLoad(target interface{}, modifiers ...tag.Modifier) {
	if err := c.Load(target, modifiers...); err != nil {
		panic(err.Error())
	}
}

func (c *envConfig) load(target interface{}) []error {
	objValue, objType, err := getIndirect(target)
	if err != nil {
		return []error{err}
	}

	errors := []error{}

	for i := 0; i < objType.NumField(); i++ {
		var (
			fieldType        = objType.Field(i)
			fieldValue       = objValue.Field(i)
			envTagValue      = fieldType.Tag.Get(envTag)
			defaultTagValue  = fieldType.Tag.Get(defaultTag)
			requiredTagValue = fieldType.Tag.Get(requiredTag)
		)

		if envTagValue == "" {
			continue
		}

		envTags := []string{
			strings.ToUpper(fmt.Sprintf("%s_%s", c.prefix, envTagValue)),
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

	return errors
}

func getExportedFields(obj interface{}) ([]*reflectField, error) {
	objValue, objType, err := getIndirect(obj)
	if err != nil {
		return nil, err
	}

	fields := []*reflectField{}
	for i := 0; i < objType.NumField(); i++ {
		field, fieldType := objValue.Field(i), objType.Field(i)
		if !isExported(fieldType.Name) {
			continue
		}

		fields = append(fields, &reflectField{
			field:     field,
			fieldType: fieldType,
		})
	}

	return fields, nil
}

func isExported(name string) bool {
	return unicode.IsUpper([]rune(name)[0])
}

func getIndirect(obj interface{}) (reflect.Value, reflect.Type, error) {
	indirect := reflect.Indirect(reflect.ValueOf(obj))
	if !indirect.IsValid() {
		return reflect.Value{}, nil, fmt.Errorf("invalid type for configuration struct")
	}

	indirectType := indirect.Type()
	if indirectType.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("invalid type for configuration struct")
	}

	return indirect, indirectType, nil
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
