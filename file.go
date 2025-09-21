package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vekio/x/fs"
	"github.com/vekio/x/fs/file"
)

// ConfigFile wraps the metadata and helpers required to manage one
// application-specific configuration file.
type ConfigFile[T Validatable] struct {
	fileManager FileManager[T]
	fileName    string
	path        string
	appName     string
	data        T
	defaultData T
}

// Validatable is implemented by configuration types that can perform their own
// validation after being loaded from disk.
type Validatable interface {
	Validate() error
}

// DirPath returns the full directory path where the application's configuration files are stored.
// It combines the configuration directory and the application's name.
func (c *ConfigFile[T]) DirPath() string {
	return filepath.Join(c.path, c.appName)
}

// Path constructs and returns the full path to the configuration file.
// It combines the directory path and the file name.
func (c *ConfigFile[T]) Path() string {
	return filepath.Join(c.DirPath(), c.fileName)
}

// Content reads and returns the content of the configuration file.
// It returns an error if the file cannot be read.
func (c *ConfigFile[T]) Content() ([]byte, error) {
	buf, err := os.ReadFile(c.Path())
	if err != nil {
		return nil, fmt.Errorf("read configuration file: %w", err)
	}
	return buf, nil
}

// Data returns the in-memory copy of the configuration that was last
// loaded from disk or passed to Init/SoftInit.
func (c *ConfigFile[T]) Data() T {
	return c.data
}

// Reload refreshes the cached configuration by pulling the latest content
// from disk using the configured file manager.
func (c *ConfigFile[T]) Reload() error {
	if err := c.fileManager.LoadDataFromFile(c.Path(), &c.data); err != nil {
		return fmt.Errorf("load configuration file: %w", err)
	}
	return nil
}

// Init initializes the configuration by ensuring that the directory and file exist,
// and by writing the initial configuration data to the file.
func (c *ConfigFile[T]) Init(data T) error {
	if err := fs.EnsureDir(c.DirPath(), fs.DefaultDirMode); err != nil {
		return fmt.Errorf("ensure config directory: %w", err)
	}
	if err := file.Touch(c.Path(), fs.DefaultFileMode); err != nil {
		return fmt.Errorf("ensure config file: %w", err)
	}
	c.data = data

	if err := c.fileManager.WriteDataToFile(c.Path(), data); err != nil {
		return fmt.Errorf("write configuration file: %w", err)
	}
	return nil
}

// SoftInit attempts to initialize the configuration by loading existing data or creating new configuration.
// It reads the configuration if the file exists or initializes it if it does not.
func (c *ConfigFile[T]) SoftInit() error {
	exists, err := file.Exists(c.Path())
	if err != nil {
		return fmt.Errorf("check configuration file: %w", err)
	}
	if !exists {
		return c.Init(c.defaultData)
	}

	if err := c.fileManager.LoadDataFromFile(c.Path(), &c.data); err != nil {
		return fmt.Errorf("load configuration file: %w", err)
	}
	return nil
}

// getFileNameForEnvironment returns a variant of the config file name that
// includes the value of APPNAME_ENV if that environment variable is defined.
func getFileNameForEnvironment(appName, configFileName string) string {
	envValue := os.Getenv(fmt.Sprintf("%s_ENV", strings.ToUpper(appName)))
	if envValue != "" {
		configFileNameSplited := strings.Split(configFileName, ".")
		if len(configFileNameSplited) > 1 {
			return fmt.Sprintf("%s.%s.%s", configFileNameSplited[0], strings.ToLower(envValue), strings.Join(configFileNameSplited[1:], "."))
		}
	}
	return configFileName
}

// validateConfigParams ensures that the minimal configuration parameters were
// provided before building a ConfigFile instance.
func validateConfigParams(path, fileName, appName string) error {
	if path == "" || fileName == "" || appName == "" {
		return fmt.Errorf("path, fileName, and appName must not be empty")
	}
	return nil
}
