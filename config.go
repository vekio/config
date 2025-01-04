package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_dir "github.com/vekio/fs/dir"
	_file "github.com/vekio/fs/file"
	"gopkg.in/yaml.v3"
)

// DefaultConfigFileName is the default file name for the application configuration.
var DefaultConfigFileName = "config.yml"

// Config manages configuration files for an application
type Config[T any] struct {
	appName  string // Name of the application
	dir      string // Directory where the configuration file is stored
	fileName string // Configuration file name
	data     *T
}

// NewConfig returns a new Config.
// It initializes the configuration with the appropriate user config directory
// and file name based on the environment.
func NewConfig[T any]() *Config[T] {
	appName := filepath.Base(os.Args[0])

	// Retrieve the user configuration directory
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("error retrieving user config directory: %v", err)
	}

	// Determine the file name based on the environment variable
	fileName := envConfigFileName(appName, DefaultConfigFileName)

	c := &Config[T]{
		dir:      dir,
		appName:  appName,
		fileName: fileName,
		data:     new(T),
	}
	return c
}

// Data returns the data [T] struct loaded.
func (c *Config[T]) Data() *T {
	return c.data
}

// AppName returns the name of the application.
func (c *Config[T]) AppName() string {
	return c.appName
}

// DirPath returns the full directory path where the application's configuration files are stored.
// It combines the configuration directory and the application's name.
func (c *Config[T]) DirPath() string {
	return filepath.Join(c.dir, c.appName)
}

// Path constructs and returns the full path to the configuration file.
// It combines the directory path and the file name.
func (c *Config[T]) Path() string {
	return filepath.Join(c.DirPath(), c.fileName)
}

// Content reads and returns the content of the configuration file.
// It returns an error if the file cannot be read.
func (c *Config[T]) Content() ([]byte, error) {
	return _file.ReadFile(c.Path())
}

// Init initializes the configuration by ensuring that the directory and file exist,
// and by writing the initial configuration data to the file.
func (c *Config[T]) Init(data T) error {
	if err := _dir.EnsureDir(c.DirPath(), _dir.DefaultDirPerms); err != nil {
		return fmt.Errorf("failed to ensure directory: %v", err)
	}
	if err := _file.CreateFile(c.Path(), _file.DefaultFilePerms); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	c.data = &data
	return c.writeDataToFile(data)
}

// SoftInit attempts to initialize the configuration by loading existing data or creating new configuration.
// It reads the configuration if the file exists or initializes it if it does not.
func (c *Config[T]) SoftInit() error {
	exists, err := _file.FileExists(c.Path())
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %v", err)
	}
	if !exists {
		return c.Init(*new(T))
	}
	return c.loadDataFromFile()
}

// writeDataToFile marshals data to YAML and writes it to the configuration file.
func (c *Config[T]) writeDataToFile(data T) error {
	buf, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	if err := _file.WriteFileContent(c.Path(), buf, _file.RestrictedFilePerms); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}
	return nil
}

// loadDataFromFile reads the configuration from the file and unmarshals it into the data structure.
func (c *Config[T]) loadDataFromFile() error {
	buf, err := os.ReadFile(c.Path())
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	if err := yaml.Unmarshal(buf, c.data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return nil
}

// envConfigFileName constructs the configuration file name based on the environment.
// It adjusts the file name to include an environment descriptor if specified.
func envConfigFileName(appName, configFileName string) string {
	envValue := os.Getenv(fmt.Sprintf("%s_ENV", strings.ToUpper(appName)))
	if envValue != "" {
		configFileNameSplited := strings.Split(configFileName, ".")
		if len(configFileNameSplited) > 1 {
			return fmt.Sprintf("%s.%s.%s", configFileNameSplited[0], strings.ToLower(envValue), strings.Join(configFileNameSplited[1:], "."))
		}
	}
	return configFileName
}
