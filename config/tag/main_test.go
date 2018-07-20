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
