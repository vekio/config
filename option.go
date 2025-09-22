package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ConfigFileOption customizes the behavior of a ConfigFile during construction.
type ConfigFileOption[T Validatable] func(*ConfigFile[T])

// WithDefault seeds the ConfigFile with a default value that will be written to
// disk during initialization when no configuration file exists yet.
func WithDefault[T Validatable](defaultData T) ConfigFileOption[T] {
	return func(c *ConfigFile[T]) {
		if c == nil {
			return
		}
		c.defaultData = defaultData
	}
}

// WithAppName overrides the application identifier used to locate the
// configuration directory and derive environment-specific overrides.
func WithAppName[T Validatable](appName string) ConfigFileOption[T] {
	return func(c *ConfigFile[T]) {
		if c == nil {
			return
		}

		trimmed := strings.TrimSpace(appName)
		if trimmed == "" {
			return
		}
		c.appName = trimmed
	}
}

// WithPath overrides the directory where configuration files are stored. When
// the provided path is empty the default user configuration directory (or the
// system temp directory as fallback) is used.
func WithPath[T Validatable](path string) ConfigFileOption[T] {
	return func(c *ConfigFile[T]) {
		if c == nil {
			return
		}

		trimmed := strings.TrimSpace(path)
		if trimmed == "" {
			return
		}
		c.path = trimmed
	}
}

// WithFilename overrides the default file name while ensuring it uses the
// extension that matches the underlying file manager.
func WithFilename[T Validatable](fileName string) ConfigFileOption[T] {
	return func(c *ConfigFile[T]) {
		if c == nil {
			return
		}

		trimmed := strings.TrimSpace(fileName)
		if trimmed == "" {
			return
		}

		base := filepath.Base(trimmed)
		if base == "" {
			return
		}

		extension := strings.TrimPrefix(c.fileManager.Extension(), ".")
		if currentExt := filepath.Ext(base); currentExt != "" {
			base = strings.TrimSuffix(base, currentExt)
		}

		c.fileName = fmt.Sprintf("%s.%s", base, extension)
	}
}
