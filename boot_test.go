package nacelle

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type BootSuite struct{}

func (s *BootSuite) TestBoot(t sweet.T) {
	ran := false
	bootstrapper := NewBootstrapper(
		"APP",
		func(processes ProcessContainer, services ServiceContainer) error {
			processes.RegisterInitializer(InitializerFunc(func(config Config) error {
				ran = true
				return nil
			}))

			return nil
		},
	)

	Expect(bootstrapper.Boot()).To(Equal(0))
	Expect(ran).To(BeTrue())
}

func (s *BootSuite) TestDefaultServices(t sweet.T) {
	serviceChecker := &struct {
		Health   Health           `service:"health"`
		Logger   Logger           `service:"logger"`
		Services ServiceContainer `service:"services"`
	}{}

	bootstrapper := NewBootstrapper(
		"APP",
		func(processes ProcessContainer, services ServiceContainer) error {
			return services.Inject(serviceChecker)
		},
	)

	Expect(bootstrapper.Boot()).To(Equal(0))
	Expect(serviceChecker.Health).NotTo(BeNil())
	Expect(serviceChecker.Logger).NotTo(BeNil())
	Expect(serviceChecker.Services).NotTo(BeNil())
}

func (s *BootSuite) TestInitFuncError(t sweet.T) {
	bootstrapper := NewBootstrapper(
		"APP",
		func(processes ProcessContainer, services ServiceContainer) error {
			return fmt.Errorf("oops")
		},
	)

	Expect(bootstrapper.Boot()).To(Equal(1))
}

func (s *BootSuite) TestLoggingInitError(t sweet.T) {
	bootstrapper := NewBootstrapper(
		"APP",
		func(processes ProcessContainer, services ServiceContainer) error {
			return nil
		},
		WithLoggingInitFunc(func(Config) (Logger, error) {
			return nil, fmt.Errorf("oops")
		}),
	)

	Expect(bootstrapper.Boot()).To(Equal(1))
}

func (s *BootSuite) TestRunnerError(t sweet.T) {
	bootstrapper := NewBootstrapper(
		"APP",
		func(processes ProcessContainer, services ServiceContainer) error {
			processes.RegisterInitializer(InitializerFunc(func(config Config) error {
				return fmt.Errorf("oops")
			}))

			return nil
		},
	)

	Expect(bootstrapper.Boot()).To(Equal(1))
}
