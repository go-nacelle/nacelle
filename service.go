package nacelle

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type (
	ServiceContainer struct {
		services map[interface{}]interface{}
	}

	ServiceInitializerFunc func(config Config, container *ServiceContainer) error
)

const (
	serviceTag  = "service"
	optionalTag = "optional"
)

var (
	ErrDuplicateServiceKey    = errors.New("duplicate service key")
	ErrUnregisteredServiceKey = errors.New("no service registered to key")
	ErrIllegalLogger          = errors.New("logger instance is not a nacelle.Logger")
)

func WrapServiceInitializerFunc(container *ServiceContainer, f ServiceInitializerFunc) InitializerFunc {
	return InitializerFunc(func(config Config) error {
		return f(config, container)
	})
}

func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		services: map[interface{}]interface{}{},
	}
}

func (c *ServiceContainer) Get(service interface{}) (interface{}, error) {
	value, ok := c.services[service]
	if !ok {
		return nil, ErrUnregisteredServiceKey
	}

	return value, nil
}

func (c *ServiceContainer) GetLogger() Logger {
	if raw, err := c.Get("logger"); err == nil {
		return raw.(Logger)

	}

	return emergencyLogger()
}

func (c *ServiceContainer) MustGet(service interface{}) interface{} {
	value, err := c.Get(service)
	if err != nil {
		panic(err.Error())
	}

	return value
}

func (c *ServiceContainer) Set(service, value interface{}) error {
	if service == "logger" {
		if _, ok := value.(Logger); !ok {
			return ErrIllegalLogger
		}
	}

	if _, ok := c.services[service]; ok {
		return ErrDuplicateServiceKey
	}

	c.services[service] = value
	return nil
}

func (c *ServiceContainer) MustSet(service, value interface{}) {
	if err := c.Set(service, value); err != nil {
		panic(err.Error())
	}
}

func (c *ServiceContainer) Inject(obj interface{}) error {
	var (
		ov = reflect.ValueOf(obj)
		oi = reflect.Indirect(ov)
		ot = oi.Type()
	)

	if oi.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < ot.NumField(); i++ {
		var (
			fieldType   = ot.Field(i)
			fieldValue  = oi.Field(i)
			serviceTag  = fieldType.Tag.Get(serviceTag)
			optionalTag = fieldType.Tag.Get(optionalTag)
		)

		if serviceTag == "" {
			continue
		}

		if err := loadServiceField(c, fieldType, fieldValue, serviceTag, optionalTag); err != nil {
			return err
		}
	}

	return nil
}

func loadServiceField(container *ServiceContainer, fieldType reflect.StructField, fieldValue reflect.Value, serviceTag, optionalTag string) error {
	if !fieldValue.IsValid() {
		return fmt.Errorf("field '%s' is invalid", fieldType.Name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field '%s' can not be set", fieldType.Name)
	}

	value, err := container.Get(serviceTag)
	if err != nil {
		if optionalTag != "" {
			val, err := strconv.ParseBool(optionalTag)
			if err != nil {
				return fmt.Errorf("field '%s' has an invalid optional tag", fieldType.Name)
			}

			if val {
				return nil
			}
		}

		return err
	}

	var (
		targetType  = fieldValue.Type()
		targetValue = reflect.ValueOf(value)
	)

	if !targetValue.Type().ConvertibleTo(targetType) {
		return fmt.Errorf("field '%s' cannot be assigned a value of type %s", fieldType.Name, targetType)
	}

	fieldValue.Set(targetValue.Convert(targetType))
	return nil
}
