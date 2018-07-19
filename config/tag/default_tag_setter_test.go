package tag

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type DefaultTagSetterSuite struct{}

func (s *DefaultTagSetterSuite) TestDefaultTagSetter(t sweet.T) {
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
