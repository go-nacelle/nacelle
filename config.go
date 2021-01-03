package nacelle

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-nacelle/config"
	"github.com/go-nacelle/log"
)

// ConfigurationRegistry is a wrapper around a bare configuration loader object.
type ConfigurationRegistry interface {
	// Register populates the given configuration object. References to the
	// object and tag modifiers may be held for later reporting.
	Register(target interface{}, modifiers ...config.TagModifier)
}

type Configurable interface {
	// RegisterConfiguration is a hook provided with a configuration registry. The
	// registry object populates configuration objects and aggregates configuration
	// and validation errors into a central location.
	//
	// This hook is called prior to the Init method of any registered initializer
	// or process.
	RegisterConfiguration(registry ConfigurationRegistry)
}

type namedInitializer interface {
	Name() string
	LogFields() log.LogFields
}

type registeredConfig struct {
	meta      namedInitializer
	target    interface{}
	loadErr   error
	modifiers []config.TagModifier
}

func newConfigurationRegistry(c config.Config, meta namedInitializer, f func(registeredConfig)) ConfigurationRegistry {
	return configurationRegisterFunc(func(target interface{}, modifiers ...config.TagModifier) {
		f(registeredConfig{
			meta:      meta,
			target:    target,
			loadErr:   c.Load(target, modifiers...),
			modifiers: modifiers,
		})
	})
}

type configurationRegisterFunc func(target interface{}, modifiers ...config.TagModifier)

func (f configurationRegisterFunc) Register(target interface{}, modifiers ...config.TagModifier) {
	f(target, modifiers...)
}

func loadConfig(processes ProcessContainer, config config.Config, logger Logger) []registeredConfig {
	logger.Info("Loading configuration")

	var configs []registeredConfig
	for i := 0; i < processes.NumInitializerPriorities(); i++ {
		for _, initializer := range processes.GetInitializersAtPriorityIndex(i) {
			if configurable, ok := initializer.Wrapped().(Configurable); ok {
				registry := newConfigurationRegistry(config, initializer, func(config registeredConfig) {
					configs = append(configs, config)
				})

				configurable.RegisterConfiguration(registry)
			}
		}
	}

	for i := 0; i < processes.NumProcessPriorities(); i++ {
		for _, process := range processes.GetProcessesAtPriorityIndex(i) {
			if configurable, ok := process.Wrapped().(Configurable); ok {
				registry := newConfigurationRegistry(config, process, func(config registeredConfig) {
					configs = append(configs, config)
				})

				configurable.RegisterConfiguration(registry)
			}
		}
	}

	return configs
}

func validateConfig(config config.Config, configs []registeredConfig, logger Logger) error {
	logger.Info("Validating configuration")

	var errors []error
	for _, c := range configs {
		logger := logger.WithFields(c.meta.LogFields())

		if c.loadErr != nil {
			logger.Error(
				"Failed to load configuration for %s (%s)",
				c.meta.Name(),
				c.loadErr.Error(),
			)

			errors = append(errors, c.loadErr)
			continue
		}

		if err := config.PostLoad(c.target); err != nil {
			logger.Error(
				"PostLoad failed for configuration target in %s (%s)",
				c.meta.Name(),
				err.Error(),
			)

			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("configuration validation failed: %d errors", len(errors))
	}

	return nil
}

func describeConfiguration(config config.Config, configs []registeredConfig, logger Logger, additionalConfigurationTargets ...interface{}) (string, error) {
	var descriptions []string
	for _, c := range additionalConfigurationTargets {
		description, err := config.Describe(c)
		if err != nil {
			return "", err
		}

		descriptions = append(descriptions, formatConfigDescription(description)...)
	}

	for _, c := range configs {
		description, err := config.Describe(c.target, c.modifiers...)
		if err != nil {
			return "", err
		}

		descriptions = append(descriptions, formatConfigDescription(description)...)
	}

	sort.Strings(descriptions)
	return strings.Join(descriptions, "\n"), nil
}

func formatConfigDescription(description *config.StructDescription) []string {
	var descriptions []string
	for _, field := range description.Fields {
		for key, value := range field.TagValues {
			var parts []string
			if field.Required {
				parts = append(parts, "required")
			}
			if field.Default != "" {
				parts = append(parts, fmt.Sprintf("default=%s", field.Default))
			}

			suffix := ""
			if len(parts) > 0 {
				suffix = fmt.Sprintf(" (%s)", strings.Join(parts, "; "))
			}

			descriptions = append(descriptions, fmt.Sprintf("(%s) %s%s", key, value, suffix))
		}
	}

	return descriptions
}
