package config

import (
	"reflect"

	"github.com/aphistic/sweet"
	"github.com/fatih/structtag"
	. "github.com/onsi/gomega"
)

type ConfigTagsSuite struct{}

func (s *ConfigTagsSuite) TestEnvTagPrefixer(t sweet.T) {
	obj, err := ApplyTagModifiers(&TempTest{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}

func (s *ConfigTagsSuite) TestDefaultTagSetter(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&TempTest{},
		NewDefaultTagSetter("X", "r"),
		NewDefaultTagSetter("Y", "null"),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"default": "r",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{
		"default": "null",
	}))
}

//
//

type TempTest struct {
	X string `env:"a" default:"q"`
	Y string
}

//
//

func gatherTags(obj interface{}, name string) map[string]string {
	var (
		ov = reflect.ValueOf(obj)
		oi = reflect.Indirect(ov)
		ot = oi.Type()
	)

	for i := 0; i < ot.NumField(); i++ {
		if fieldType := ot.Field(i); fieldType.Name == name {
			if tags, ok := getTags(fieldType); ok {
				return decomposeTags(tags)
			}
		}
	}

	return nil
}

func decomposeTags(tags *structtag.Tags) map[string]string {
	fieldTags := map[string]string{}

	for _, name := range tags.Keys() {
		tag, _ := tags.Get(name)
		fieldTags[name] = tag.Name
	}

	return fieldTags
}
