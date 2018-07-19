package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/efritz/nacelle/config/tag"
	"github.com/efritz/nacelle/logging"
)

// LoggingConfig decorates another config loader object with logging.
// After each successful load, the populated configuration object is
// serialized as fields and output at the info level.
type LoggingConfig struct {
	Config
	logger logging.Logger
}

const (
	maskTag    = "mask"
	displayTag = "display"
)

func NewLoggingConfig(config Config, logger logging.Logger) Config {
	return &LoggingConfig{
		Config: config,
		logger: logger,
	}
}

func (c *LoggingConfig) Load(target interface{}, modifiers ...tag.TagModifier) error {
	if err := c.Config.Load(target); err != nil {
		return err
	}

	m, err := dumpChunk(target)
	if err != nil {
		return fmt.Errorf("failed to serialize config (%s)", err.Error())
	}

	c.logger.InfoWithFields(m, "Config loaded from environment")
	return nil
}

func (c *LoggingConfig) MustLoad(target interface{}, modifiers ...tag.TagModifier) {
	if err := c.Load(target, modifiers...); err != nil {
		panic(err.Error())
	}
}

func dumpChunk(obj interface{}) (map[string]interface{}, error) {
	objValue, objType, err := getIndirect(obj)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	for i := 0; i < objType.NumField(); i++ {
		var (
			fieldType  = objType.Field(i)
			fieldValue = objValue.Field(i)
			// envTagValue     = fieldType.Tag.Get(envTag)
			maskTagValue    = fieldType.Tag.Get(maskTag)
			displayTagValue = fieldType.Tag.Get(displayTag)
			displayName     = fieldType.Name
		)

		if displayTagValue != "" {
			displayName = displayTagValue
		}

		if maskTagValue != "" {
			val, err := strconv.ParseBool(maskTagValue)
			if err != nil {
				return nil, fmt.Errorf("field '%s' has an invalid mask tag", fieldType.Name)
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
				return nil, err
			}

			m[displayName] = string(data)
		}
	}

	return m, nil
}
