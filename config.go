package config

import (
	"fmt"
	"strings"
)

// NewYAMLConfigFile constructs a ConfigFile that persists data as YAML inside
// the default configuration directory. The base path can be overridden with
// WithPath.
func NewYAMLConfigFile[T Validatable](options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	return newConfigFile(NewYAMLFileManager[T](), options...)
}

// NewJSONConfigFile constructs a ConfigFile that persists data as JSON inside
// the default configuration directory. The base path can be overridden with
// WithPath.
func NewJSONConfigFile[T Validatable](options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	return newConfigFile(NewJSONFileManager[T](), options...)
}

// NewDefaultConfigFile builds a YAML configuration backed by the user's
// configuration directory (falling back to the OS temp dir when unavailable).
func NewDefaultConfigFile[T Validatable](options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	return NewYAMLConfigFile(options...)
}

func newConfigFile[T Validatable](manager FileManager[T], options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	if manager == nil {
		return nil, fmt.Errorf("config: file manager must not be nil")
	}

	extension := strings.TrimPrefix(manager.Extension(), ".")

	fileName := "config"
	if extension != "" {
		fileName = fmt.Sprintf("%s.%s", fileName, extension)
	}

	c := &ConfigFile[T]{
		fileManager: manager,
		appName:     defaultAppName(),
		path:        defaultConfigPath(),
		fileName:    fileName,
		defaultData: *new(T),
	}

	for _, option := range options {
		if option != nil {
			option(c)
		}
	}

	return c, nil
}
