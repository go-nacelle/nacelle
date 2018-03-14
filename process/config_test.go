package process

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestHTTPCertConfiguration(t sweet.T) {
	// Non-TLS
	c := &HTTPConfig{}
	Expect(c.PostLoad()).To(BeNil())

	// Successful TLS Config
	c = &HTTPConfig{HTTPCertFile: "cert", HTTPKeyFile: "key"}
	Expect(c.PostLoad()).To(BeNil())

	// Incomplete
	c = &HTTPConfig{HTTPCertFile: "cert"}
	Expect(c.PostLoad()).To(Equal(ErrBadCertConfig))

	c = &HTTPConfig{HTTPKeyFile: "key"}
	Expect(c.PostLoad()).To(Equal(ErrBadCertConfig))
}
