package nacelle

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle/log"
	. "github.com/onsi/gomega"
)

type ServiceSuite struct{}

func (s *ServiceSuite) TestGetLogger(t sweet.T) {
	container, err := MakeServiceContainer()
	Expect(err).To(BeNil())

	logger, _ := log.InitGomolShim(&LoggingConfig{})
	Expect(container.Set("logger", logger)).To(BeNil())
	Expect(container.GetLogger()).To(Equal(logger))
}

func (s *ServiceSuite) TestGetUnregisteredLogger(t sweet.T) {
	container, err := MakeServiceContainer()
	Expect(err).To(BeNil())
	Expect(container.GetLogger()).NotTo(BeNil())
}

func (s *ServiceSuite) TestSetBadLogger(t sweet.T) {
	container, err := MakeServiceContainer()
	Expect(err).To(BeNil())
	Expect(container.Set("logger", struct{}{})).To(MatchError("logger instance is not a nacelle.Logger"))
}
