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
	fileName := getFileNameForEnvironment(c.DirPath(), c.appName, c.fileName)
	return filepath.Join(c.DirPath(), fileName)
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

// getFileNameForEnvironment
func getFileNameForEnvironment(dirPath, appName, configFileName string) string {
	envVarName := fmt.Sprintf("%s_ENV", strings.ToUpper(appName))
	envValue := strings.TrimSpace(os.Getenv(envVarName))
	if envValue == "" || strings.EqualFold(envValue, "pro") {
		return configFileName
	}

	base := filepath.Base(configFileName)
	extension := filepath.Ext(configFileName)
	lowerEnv := strings.ToLower(envValue)
	candidateFileName := fmt.Sprintf("%s.%s.%s", base, lowerEnv, extension)

	candidatePath := filepath.Join(dirPath, candidateFileName)
	exists, err := file.Exists(candidatePath)
	if err != nil {
		return configFileName
	}
	if exists {
		return candidateFileName
	}

	return configFileName
}
