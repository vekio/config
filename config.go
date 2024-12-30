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

type Validatable interface {
	Validate() error
}

// Config manages configuration files for an application.
type Config struct {
	appName string // Name of the application
	dir     string // Directory where the config files are stored
	file    string // Configuration file name
}

func newConfig() *Config {
	// Extract the executable name from the first argument
	appName := filepath.Base(os.Args[0])

	// User config directory
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("failed to get user config directory: %s", err)
	}

	// Check for the app environment variable
	envValue := os.Getenv(fmt.Sprintf("%s_ENV", strings.ToUpper(appName)))
	configFile := "config.yml"

	if envValue != "" && envValue == "develop" {
		configFile = "config.dev.yml"
	}

	c := &Config{
		dir:     dir,
		appName: appName,
		file:    configFile,
	}
	return c
}

func (c *Config) AppName() string {
	return c.appName
}

func (c *Config) DirPath() string {
	return filepath.Join(c.dir, c.appName)
}

func (c *Config) Path() string {
	return filepath.Join(c.DirPath(), c.file)
}

func (c *Config) Content() ([]byte, error) {
	return os.ReadFile(c.Path())
}

func (c *Config) Init() error {
	err := _dir.EnsureDir(c.DirPath(), _dir.DefaultDirPerms)
	if err != nil {
		return err
	}
	if err := _file.CreateFile(c.Path(), _file.DefaultFilePerms); err != nil {
		return err
	}
	// TODO writing default configuration as YAML
	// defaultConfig := new(T) // Create a zero value for T to marshal into YAML
	// data, err := yaml.Marshal(defaultConfig)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal default config: %w", err)
	// }
	// _, err = file.Write(data)
	// if err != nil {
	// 	return fmt.Errorf("failed to write default config to file %s: %w", cm.Path(), err)
	// }
	return nil
}

func (c *Config) SoftInit() error {
	exists, err := _file.FileExists(c.Path())
	if err != nil {
		return err
	}
	if !exists {
		return c.Init()
	}
	return nil
}

func (c *Config) Load(data Validatable) error {
	buf, err := c.Content()
	if err != nil {
		return err
	}
	// Deserialize the configuration file.
	if err := yaml.Unmarshal(buf, data); err != nil {
		return err
	}
	// Validate configuration data.
	if err := data.Validate(); err != nil {
		return err
	}
	return nil
}
