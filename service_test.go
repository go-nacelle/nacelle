package nacelle

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle/log"
	. "github.com/onsi/gomega"
)

type ServiceSuite struct{}

func (s *ServiceSuite) TestGetLogger(t sweet.T) {
	container := NewServiceContainer()
	logger, _ := log.InitGomolShim(&LoggingConfig{})
	err := container.Set("logger", logger)
	Expect(err).To(BeNil())
	Expect(container.GetLogger()).To(Equal(logger))
}

func (s *ServiceSuite) TestGetUnregisteredLogger(t sweet.T) {
	Expect(NewServiceContainer().GetLogger()).NotTo(BeNil())
}

func (s *ServiceSuite) TestSetBadLogger(t sweet.T) {
	Expect(NewServiceContainer().Set("logger", struct{}{})).To(MatchError("logger instance is not a nacelle.Logger"))
}
