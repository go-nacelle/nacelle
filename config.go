package nacelle

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-nacelle/config/v2"
	"github.com/go-nacelle/process/v2"
)

// ConfigurationRegistry is a wrapper around a bare configuration loader object.
type ConfigurationRegistry interface {
	// Register populates the given configuration object. References to the
	// object and tag modifiers may be held for later reporting.
	Register(target interface{}, modifiers ...TagModifier)
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

type registeredConfig struct {
	meta      *process.Meta
	target    interface{}
	loadErr   error
	modifiers []TagModifier
}

func newConfigurationRegistry(c *Config, meta *process.Meta, f func(registeredConfig)) ConfigurationRegistry {
	return configurationRegisterFunc(func(target interface{}, modifiers ...TagModifier) {
		f(registeredConfig{
			meta:      meta,
			target:    target,
			loadErr:   c.Load(target, modifiers...),
			modifiers: modifiers,
		})
	})
}

type configurationRegisterFunc func(target interface{}, modifiers ...TagModifier)

func (f configurationRegisterFunc) Register(target interface{}, modifiers ...TagModifier) {
	f(target, modifiers...)
}

func loadConfig(processes *ProcessContainer, config *Config, logger Logger) []registeredConfig {
	logger.Info("Loading configuration")

	var configs []registeredConfig
	for _, meta := range processes.Meta() {
		if configurable, ok := meta.Wrapped().(Configurable); ok {
			registry := newConfigurationRegistry(config, meta, func(config registeredConfig) {
				configs = append(configs, config)
			})

			configurable.RegisterConfiguration(registry)
		}
	}

	return configs
}

func validateConfig(config *Config, configs []registeredConfig, logger Logger) error {
	logger.Info("Validating configuration")

	var errors []error
	for _, c := range configs {
		logger := logger.WithFields(LogFields(c.meta.Metadata()))

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

func describeConfiguration(config *Config, configs []registeredConfig, logger Logger, additionalConfigurationTargets ...interface{}) (string, error) {
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
