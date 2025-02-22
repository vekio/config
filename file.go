package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_dir "github.com/vekio/fs/dir"
	_file "github.com/vekio/fs/file"
)

type ConfigFile[T Validatable] struct {
	fileManager FileManager[T]
	fileName    string
	path        string
	appName     string
	data        T
	defaultData T
}

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
	return _file.ReadFile(c.Path())
}

func (c *ConfigFile[T]) Data() T {
	return c.data
}

// Init initializes the configuration by ensuring that the directory and file exist,
// and by writing the initial configuration data to the file.
func (c *ConfigFile[T]) Init(data T) error {
	if err := _dir.EnsureDir(c.DirPath(), _dir.DefaultDirPerms); err != nil {
		return err
	}
	if err := _file.CreateFile(c.Path(), _file.DefaultFilePerms); err != nil {
		return err
	}
	c.data = data
	return c.fileManager.WriteDataToFile(c.Path(), data)
}

// SoftInit attempts to initialize the configuration by loading existing data or creating new configuration.
// It reads the configuration if the file exists or initializes it if it does not.
func (c *ConfigFile[T]) SoftInit() error {
	exists, err := _file.FileExists(c.Path())
	if err != nil {
		return err
	}
	if !exists {
		return c.Init(c.defaultData)
	}

	err = c.fileManager.LoadDataFromFile(c.Path(), &c.data)
	if err != nil {
		return err
	}
	return nil
}

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

func validateConfigParams(path, fileName, appName string) error {
	if path == "" || fileName == "" || appName == "" {
		return fmt.Errorf("path, fileName, and appName must not be empty")
	}
	return nil
}
