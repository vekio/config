package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewYAMLConfigFile creates a typed YAML configuration file inside the user's
// configuration directory.
func NewYAMLConfigFile[T Validatable](appName, fileName string) (*ConfigFile[T], error) {
	return newConfigFile(appName, fileName, yamlCodec[T]{})
}

// NewJSONConfigFile creates a typed JSON configuration file inside the user's
// configuration directory.
func NewJSONConfigFile[T Validatable](appName, fileName string) (*ConfigFile[T], error) {
	return newConfigFile(appName, fileName, jsonCodec[T]{})
}

func newConfigFile[T Validatable](appName, fileName string, valueCodec codec[T]) (*ConfigFile[T], error) {
	cfg := &ConfigFile[T]{
		codec:    valueCodec,
		fileName: fileName,
		appName:  appName,
	}
	if err := validateConfigFile(cfg); err != nil {
		return nil, err
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
	cfg.baseDir = baseDir
	return cfg, nil
}

func validateConfigFile[T Validatable](cfg *ConfigFile[T]) error {
	if !isPathElement(cfg.appName) {
		return fmt.Errorf("application name %q must be a single path element", cfg.appName)
	}
	if !isPathElement(cfg.fileName) {
		return fmt.Errorf("configuration filename %q must be a single path element", cfg.fileName)
	}
	return nil
}

func isPathElement(value string) bool {
	return value != "" && value == strings.TrimSpace(value) && value != "." && value != ".." &&
		filepath.Base(value) == value && !strings.Contains(value, "\\") && !strings.ContainsRune(value, 0)
}
