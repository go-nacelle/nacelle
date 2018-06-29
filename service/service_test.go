package service

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/nacelle/logging"
)

type ServiceSuite struct{}

func (s *ServiceSuite) TestGetLogger(t sweet.T) {
	container, err := NewContainer()
	Expect(err).To(BeNil())

	logger, err := logging.InitGomolShim(&logging.Config{
		LogLevel: "warning",
	})

	Expect(err).To(BeNil())
	Expect(container.Set("logger", logger)).To(BeNil())
	Expect(container.GetLogger()).To(Equal(logger))
}

func (s *ServiceSuite) TestGetUnregisteredLogger(t sweet.T) {
	container, err := NewContainer()
	Expect(err).To(BeNil())
	Expect(container.GetLogger()).NotTo(BeNil())
}

func (s *ServiceSuite) TestSetBadLogger(t sweet.T) {
	container, err := NewContainer()
	Expect(err).To(BeNil())
	Expect(container.Set("logger", struct{}{})).To(MatchError("logger instance is not a nacelle.Logger"))
}
