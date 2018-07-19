package tag

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type EnvTagPrefixerSuite struct{}

func (s *EnvTagPrefixerSuite) TestEnvTagPrefixer(t sweet.T) {
	obj, err := ApplyTagModifiers(&TempTest{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}
