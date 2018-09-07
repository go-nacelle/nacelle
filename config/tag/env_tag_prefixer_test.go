package tag

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type EnvTagPrefixerSuite struct{}

func (s *EnvTagPrefixerSuite) TestEnvTagPrefixer(t sweet.T) {
	obj, err := ApplyModifiers(&BasicConfig{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}

func (s *EnvTagPrefixerSuite) TestEnvTagPrefixerEmbedded(t sweet.T) {
	obj, err := ApplyModifiers(&ParentConfig{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}
