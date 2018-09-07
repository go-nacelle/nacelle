package tag

import (
	"reflect"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	"github.com/fatih/structtag"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&EnvTagPrefixerSuite{})
		s.AddSuite(&DefaultTagSetterSuite{})
	})
}

//
//

func gatherTags(obj interface{}, name string) map[string]string {
	var (
		objValue = reflect.Indirect(reflect.ValueOf(obj))
		objType  = objValue.Type()
	)

	return gatherTagsStruct(objValue, objType, name)
}

func gatherTagsStruct(objValue reflect.Value, objType reflect.Type, name string) map[string]string {
	for i := 0; i < objType.NumField(); i++ {
		var (
			field     = objValue.Field(i)
			fieldType = objType.Field(i)
		)

		if fieldType.Anonymous {
			if tags := gatherTagsStruct(field, fieldType.Type, name); tags != nil {
				return tags
			}
		}

		if fieldType.Name == name {
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
