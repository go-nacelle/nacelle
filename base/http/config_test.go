package http

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestHTTPCertConfiguration(t sweet.T) {
	// Non-TLS
	c := &Config{}
	Expect(c.PostLoad()).To(BeNil())

	// Successful TLS Config
	c = &Config{HTTPCertFile: "cert", HTTPKeyFile: "key"}
	Expect(c.PostLoad()).To(BeNil())

	// Incomplete
	c = &Config{HTTPCertFile: "cert"}
	Expect(c.PostLoad()).To(Equal(ErrBadCertConfig))

	c = &Config{HTTPKeyFile: "key"}
	Expect(c.PostLoad()).To(Equal(ErrBadCertConfig))
}
