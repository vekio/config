package config

import "os"

var DefaultConfigFileName string = "config.yml"

func NewYAMLConfigFile[T Validatable](path, fileName, appName string, options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	if err := validateConfigParams(path, fileName, appName); err != nil {
		return nil, err
	}

	finalFileName := getFileNameForEnvironment(appName, fileName)
	c := &ConfigFile[T]{
		fileManager: NewYAMLFileManager[T](),
		fileName:    finalFileName,
		path:        path,
		appName:     appName,
		defaultData: *new(T),
	}

	for _, option := range options {
		option(c)
	}
	return c, nil
}

func NewJSONConfigFile[T Validatable](path, fileName, appName string, options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	if err := validateConfigParams(path, fileName, appName); err != nil {
		return nil, err
	}

	finalFileName := getFileNameForEnvironment(appName, fileName)
	c := &ConfigFile[T]{
		fileManager: NewJSONFileManager[T](),
		fileName:    finalFileName,
		path:        path,
		appName:     appName,
		defaultData: *new(T),
	}

	for _, option := range options {
		option(c)
	}
	return c, nil
}

func NewDefaultConfigFile[T Validatable](appName string, options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		// return nil, fmt.Errorf("error retrieving user config directory: %w", err)
		dir = os.TempDir()
	}
	return NewYAMLConfigFile[T](dir, DefaultConfigFileName, appName, options...)
}
