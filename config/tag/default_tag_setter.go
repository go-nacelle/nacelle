package tag

import (
	"reflect"

	"github.com/fatih/structtag"
)

// DefaultTagSetter is a tag modifier which sets the value of the default
// tag for a particular field. This is used to change the default values
// provided by third party libraries (for which a source change would be
// otherwise required).
type DefaultTagSetter struct {
	field        string
	defaultValue string
}

const (
	defaultTag = "default"
)

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
