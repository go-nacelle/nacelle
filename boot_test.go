package nacelle

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-nacelle/config/v3"
	"github.com/go-nacelle/log/v2"
	"github.com/go-nacelle/process/v2"
	"github.com/go-nacelle/service/v2"
	"github.com/stretchr/testify/assert"
)

func TestBoot(t *testing.T) {
	ran := false
	bootstrapper := NewBootstrapper(
		"APP",
		func(ctx context.Context, processes *ProcessContainerBuilder, services *ServiceContainer) error {
			processes.RegisterInitializer(InitializerFunc(func(ctx context.Context) error {
				ran = true
				return nil
			}))

			return nil
		},
	)

	assert.Equal(t, 0, bootstrapper.Boot())
	assert.True(t, ran)
}

func TestDefaultServices(t *testing.T) {
	serviceChecker := &struct {
		Health   *Health           `service:"health"`
		Logger   Logger            `service:"logger"`
		Services *ServiceContainer `service:"services"`
	}{}

	bootstrapper := NewBootstrapper(
		"APP",
		func(ctx context.Context, processes *ProcessContainerBuilder, services *ServiceContainer) error {

			return service.Inject(ctx, services, serviceChecker)
		},
	)

	assert.Equal(t, 0, bootstrapper.Boot())
	assert.NotNil(t, serviceChecker.Health)
	assert.NotNil(t, serviceChecker.Logger)
	assert.NotNil(t, serviceChecker.Services)
}

func TestDefaultServicesInContext(t *testing.T) {
	bootstrapper := NewBootstrapper(
		"APP",
		func(ctx context.Context, processes *ProcessContainerBuilder, services *ServiceContainer) error {
			assert.NotNil(t, process.HealthFromContext(ctx))
			assert.NotEqual(t, log.NewNilLogger(), log.FromContext(ctx))
			assert.NotNil(t, service.FromContext(ctx))
			assert.NotNil(t, config.FromContext(ctx))

			return nil
		},
		WithContextFilter(func(ctx context.Context) context.Context {
			assert.NotNil(t, process.HealthFromContext(ctx))
			assert.NotEqual(t, log.NewNilLogger(), log.FromContext(ctx))
			assert.NotNil(t, service.FromContext(ctx))
			assert.NotNil(t, config.FromContext(ctx))

			return ctx
		}),
	)

	assert.Equal(t, 0, bootstrapper.Boot())
}

func TestInitFuncError(t *testing.T) {
	bootstrapper := NewBootstrapper(
		"APP",
		func(ctx context.Context, processes *ProcessContainerBuilder, services *ServiceContainer) error {
			return fmt.Errorf("oops")
		},
	)

	assert.Equal(t, 1, bootstrapper.Boot())
}

func TestLoggingInitError(t *testing.T) {
	bootstrapper := NewBootstrapper(
		"APP",
		func(ctx context.Context, processes *ProcessContainerBuilder, services *ServiceContainer) error {
			return nil
		},
		WithLoggingInitFunc(func(*Config) (Logger, error) {
			return nil, fmt.Errorf("oops")
		}),
	)

	assert.Equal(t, 1, bootstrapper.Boot())
}

func TestRunnerError(t *testing.T) {
	bootstrapper := NewBootstrapper(
		"APP",
		func(ctx context.Context, processes *ProcessContainerBuilder, services *ServiceContainer) error {
			processes.RegisterInitializer(InitializerFunc(func(ctx context.Context) error {
				return fmt.Errorf("oops")
			}))

			return nil
		},
	)

	assert.Equal(t, 1, bootstrapper.Boot())
}
