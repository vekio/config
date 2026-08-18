package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

const defaultYAMLFileName = "config.yml"

// NewDefaultConfigFile creates a typed YAML configuration file inside the
// user's configuration directory.
func NewDefaultConfigFile[T Validatable](appName string) (*ConfigFile[T], error) {
	baseDir, err := defaultConfigDir()
	if err != nil {
		return nil, err
	}
	return NewYAMLConfigFile[T](baseDir, appName, defaultYAMLFileName)
}

// NewYAMLConfigFile creates a typed YAML configuration file.
func NewYAMLConfigFile[T Validatable](baseDir, appName, fileName string) (*ConfigFile[T], error) {
	return newConfigFile(baseDir, appName, fileName, yamlCodec[T]{})
}

// NewJSONConfigFile creates a typed JSON configuration file.
func NewJSONConfigFile[T Validatable](baseDir, appName, fileName string) (*ConfigFile[T], error) {
	return newConfigFile(baseDir, appName, fileName, jsonCodec[T]{})
}

func newConfigFile[T Validatable](baseDir, appName, fileName string, valueCodec codec[T]) (*ConfigFile[T], error) {
	baseDir, err := cleanPath(baseDir)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration directory: %w", err)
	}
	cfg := &ConfigFile[T]{
		codec:    valueCodec,
		fileName: fileName,
		baseDir:  baseDir,
		appName:  appName,
	}
	if err := validateConfigFile(cfg); err != nil {
		return nil, err
	}
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
