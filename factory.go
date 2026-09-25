package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewYAMLConfigFile creates a typed YAML configuration file inside the user's
// configuration directory.
func NewYAMLConfigFile[T Validatable](appName, fileName string, options ...Option[T]) (*ConfigFile[T], error) {
	return newConfigFile(appName, fileName, yamlCodec[T]{}, options...)
}

// NewJSONConfigFile creates a typed JSON configuration file inside the user's
// configuration directory.
func NewJSONConfigFile[T Validatable](appName, fileName string, options ...Option[T]) (*ConfigFile[T], error) {
	return newConfigFile(appName, fileName, jsonCodec[T]{}, options...)
}

func newConfigFile[T Validatable](appName, fileName string, valueCodec codec[T], options ...Option[T]) (*ConfigFile[T], error) {
	cfg := &ConfigFile[T]{
		codec:   valueCodec,
		appName: appName,
	}
	if err := validateConfigFile(appName, fileName); err != nil {
		return nil, err
	}
	for _, option := range options {
		if option == nil {
			return nil, fmt.Errorf("configuration file option must not be nil")
		}
		option(cfg)
	}
	if path := os.Getenv(cfg.PathEnvVar()); strings.TrimSpace(path) != "" {
		if err := cfg.SetPath(path); err != nil {
			return nil, fmt.Errorf("invalid %s: %w", cfg.PathEnvVar(), err)
		}
		return cfg, nil
	}

	baseDir, err := defaultConfigDir()
	if err != nil {
		return nil, err
	}
	cfg.path = filepath.Join(baseDir, appName, fileName)
	return cfg, nil
}

func validateConfigFile(appName, fileName string) error {
	if !isPathElement(appName) {
		return fmt.Errorf("application name %q must be a single path element", appName)
	}
	if !isPathElement(fileName) {
		return fmt.Errorf("configuration filename %q must be a single path element", fileName)
	}
	return nil
}

func isPathElement(value string) bool {
	return value != "" && value == strings.TrimSpace(value) && value != "." && value != ".." &&
		filepath.Base(value) == value && !strings.Contains(value, "\\") && !strings.ContainsRune(value, 0)
}
