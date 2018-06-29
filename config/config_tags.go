package config

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/fatih/structtag"
)

type (
	// TagModifier is an interface that rewrites a set of tags for a struct
	// field. This interface is used by the ApplyTagModifiers function.
	TagModifier interface {
		// AlterFieldTag modifies the tags reference by setting or deleting
		// values from the given tag wrapper. The given tag wrapper is the
		// parsed version of the tag for the given field. Returns an error
		// if there is an internal consistency problem.
		AlterFieldTag(field reflect.StructField, tags *structtag.Tags) error
	}

	// EnvTagPrefixer is a tag modifier which adds a prefix to the values of
	// `env` tags. This can be used to register one config multiple times and
	// have their initialization be read from different environment variables.
	EnvTagPrefixer struct {
		prefix string
	}

	// DefaultTagSetter is a tag modifier which sets the value of the default
	// tag for a particular field. This is used to change the default values
	// provided by third party libraries (for which a source change would be
	// otherwise required).
	DefaultTagSetter struct {
		field        string
		defaultValue string
	}
)

// NewEnvTagPrefixer creates a new EnvTagPrefixer.
func NewEnvTagPrefixer(prefix string) TagModifier {
	return &EnvTagPrefixer{
		prefix: prefix,
	}
}

// AlterFieldTag adds the env prefixer's prefix to the `env` tag value, if one is set.
func (p *EnvTagPrefixer) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tag, err := tags.Get(envTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  envTag,
		Name: fmt.Sprintf("%s_%s", p.prefix, tag.Name),
	})
}

// NewDefaultTagSetter creates a new DefaultTagSetter.
func NewDefaultTagSetter(field string, defaultValue string) TagModifier {
	return &DefaultTagSetter{
		field:        field,
		defaultValue: defaultValue,
	}
}

// AlterFieldTag sets the value of the default tag if the field matches the target name.
func (s *DefaultTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	if fieldType.Name != s.field {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  defaultTag,
		Name: s.defaultValue,
	})
}

// ApplyTagModifiers returns a new struct with a dynamic type whose fields
// are equivalent to the given object but whose field tags are run through
// each tag modifier in sequence.
func ApplyTagModifiers(obj interface{}, modifiers ...TagModifier) (modified interface{}, err error) {
	modified = obj
	for _, modifier := range modifiers {
		modified, err = apply(modified, modifier)
		if err != nil {
			return
		}
	}

	return
}

// MustApplyTagModifiers calls ApplyTagModifiers and panics on error.
func MustApplyTagModifiers(obj interface{}, modifiers ...TagModifier) interface{} {
	modified, err := ApplyTagModifiers(obj, modifiers...)
	if err != nil {
		panic(err.Error())
	}

	return modified
}

func apply(p interface{}, modifier TagModifier) (interface{}, error) {
	val := reflect.ValueOf(p)
	newType, err := makeType(val.Type().Elem(), modifier)
	if err != nil {
		return nil, err
	}

	return reflect.NewAt(newType, unsafe.Pointer(val.Pointer())).Interface(), nil
}

func makeType(t reflect.Type, modifier TagModifier) (reflect.Type, error) {
	switch t.Kind() {
	case reflect.Struct:
		return makeStructType(t, modifier)

	case reflect.Ptr:
		inner, err := makeType(t.Elem(), modifier)
		if err != nil {
			return nil, err
		}

		return reflect.PtrTo(inner), nil

	case reflect.Array:
		inner, err := makeType(t.Elem(), modifier)
		if err != nil {
			return nil, err
		}

		return reflect.ArrayOf(t.Len(), inner), nil

	case reflect.Slice:
		inner, err := makeType(t.Elem(), modifier)
		if err != nil {
			return nil, err
		}

		return reflect.SliceOf(inner), nil

	case reflect.Map:
		key, err := makeType(t.Key(), modifier)
		if err != nil {
			return nil, err
		}

		val, err := makeType(t.Elem(), modifier)
		if err != nil {
			return nil, err
		}

		return reflect.MapOf(key, val), nil

	case
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer,
		reflect.Interface:
		return nil, fmt.Errorf("unsupported type %s", t.Kind().String())

	default:
		return t, nil
	}
}

func makeStructType(structType reflect.Type, modifier TagModifier) (reflect.Type, error) {
	fields := []reflect.StructField{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if !isExported(field.Name) {
			continue
		}

		fieldType, err := makeType(field.Type, modifier)
		if err != nil {
			return nil, err
		}

		field.Type = fieldType

		if tags, ok := getTags(field); ok {
			if err := modifier.AlterFieldTag(field, tags); err != nil {
				return nil, err
			}

			field.Tag = reflect.StructTag(tags.String())
		}

		fields = append(fields, field)
	}

	return reflect.StructOf(fields), nil
}

func getTags(field reflect.StructField) (*structtag.Tags, bool) {
	tags, err := structtag.Parse(string(field.Tag))
	if err == nil {
		return tags, true
	}

	return nil, false
}
