package config

import (
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ConfigSuite{})
		s.AddSuite(&EnvSourcerSuite{})
		s.AddSuite(&JSONSuite{})
		s.AddSuite(&LoggingConfigSuite{})
		s.AddSuite(&MultiSourcerSuite{})
		s.AddSuite(&FileSourcerSuite{})
	})
}

//
//

func ensureEquals(sourcer Sourcer, values []string, expected string) {
	val, _, ok := sourcer.Get(values)
	Expect(ok).To(BeTrue())
	Expect(val).To(Equal(expected))
}

func ensureMatches(sourcer Sourcer, values []string, expected string) {
	val, _, ok := sourcer.Get(values)
	Expect(ok).To(BeTrue())
	Expect(val).To(MatchJSON(expected))
}

func ensureMissing(sourcer Sourcer, values []string) {
	_, _, ok := sourcer.Get(values)
	Expect(ok).To(BeFalse())
}
