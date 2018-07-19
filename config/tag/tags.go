package tag

import (
	"fmt"
	"reflect"
	"unicode"
	"unsafe"

	"github.com/fatih/structtag"
)

// Modifier is an interface that rewrites a set of tags for a struct
// field. This interface is used by the ApplyModifiers function.
type Modifier interface {
	// AlterFieldTag modifies the tags reference by setting or deleting
	// values from the given tag wrapper. The given tag wrapper is the
	// parsed version of the tag for the given field. Returns an error
	// if there is an internal consistency problem.
	AlterFieldTag(field reflect.StructField, tags *structtag.Tags) error
}

// ApplyModifiers returns a new struct with a dynamic type whose fields
// are equivalent to the given object but whose field tags are run through
// each tag modifier in sequence.
func ApplyModifiers(
	obj interface{},
	modifiers ...Modifier,
) (
	modified interface{},
	err error,
) {
	modified = obj
	for _, modifier := range modifiers {
		modified, err = apply(modified, modifier)
		if err != nil {
			return
		}
	}

	return
}

func apply(p interface{}, modifier Modifier) (interface{}, error) {
	val := reflect.ValueOf(p)
	newType, err := makeType(val.Type().Elem(), modifier)
	if err != nil {
		return nil, err
	}

	return reflect.NewAt(newType, unsafe.Pointer(val.Pointer())).Interface(), nil
}

func makeType(t reflect.Type, modifier Modifier) (reflect.Type, error) {
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

func makeStructType(structType reflect.Type, modifier Modifier) (reflect.Type, error) {
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

func isExported(name string) bool {
	return unicode.IsUpper([]rune(name)[0])
}

func getTags(field reflect.StructField) (*structtag.Tags, bool) {
	tags, err := structtag.Parse(string(field.Tag))
	if err == nil {
		return tags, true
	}

	return nil, false
}
