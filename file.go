package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	xfile "github.com/vekio/x/file"
)

const (
	defaultDirMode  = 0o700
	defaultFileMode = 0o600
)

// ConfigFile wraps the metadata and helpers required to manage one
// application-specific configuration file.
type ConfigFile[T Validatable] struct {
	codec        codec[T]
	fileName     string
	baseDir      string
	appName      string
	pathOverride string
}

// Validatable is implemented by configuration types that can perform their own
// validation after being loaded from disk.
type Validatable interface {
	Validate() error
}

// Path constructs and returns the full path to the configuration file.
// It combines the base directory, application name, and file name.
func (c *ConfigFile[T]) Path() string {
	if c.pathOverride != "" {
		return filepath.Clean(c.pathOverride)
	}
	return filepath.Join(c.baseDir, c.appName, c.fileName)
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

// Load reads and validates the configuration from disk.
func (c *ConfigFile[T]) Load() (T, error) {
	var data T
	content, err := os.ReadFile(c.Path())
	if err != nil {
		return data, fmt.Errorf("read configuration file: %w", err)
	}
	if err := c.codec.unmarshal(content, &data); err != nil {
		return data, fmt.Errorf("decode configuration file: %w", err)
	}
	if err := data.Validate(); err != nil {
		return data, fmt.Errorf("validate configuration: %w", err)
	}
	return data, nil
}

// Save validates and writes the configuration to disk, replacing an existing
// file atomically.
func (c *ConfigFile[T]) Save(data T) error {
	content, err := c.encode(data)
	if err != nil {
		return err
	}
	if err := xfile.EnsureParentDir(c.Path(), defaultDirMode); err != nil {
		return fmt.Errorf("ensure configuration directory: %w", err)
	}
	if err := xfile.WriteAtomic(c.Path(), content, defaultFileMode); err != nil {
		return fmt.Errorf("write configuration file: %w", err)
	}
	return nil
}

// Create validates and writes a new configuration file. It returns an error
// wrapping os.ErrExist when the file already exists.
func (c *ConfigFile[T]) Create(data T) error {
	content, err := c.encode(data)
	if err != nil {
		return err
	}
	if err := xfile.EnsureParentDir(c.Path(), defaultDirMode); err != nil {
		return fmt.Errorf("ensure configuration directory: %w", err)
	}
	if err := xfile.WriteExclusive(c.Path(), content, defaultFileMode); err != nil {
		return fmt.Errorf("create configuration file: %w", err)
	}
	return nil
}

func (c *ConfigFile[T]) encode(data T) ([]byte, error) {
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}
	content, err := c.codec.marshal(data)
	if err != nil {
		return nil, fmt.Errorf("encode configuration: %w", err)
	}
	return content, nil
}

// LoadOrCreate loads an existing configuration or saves and returns defaultData
// when the file does not exist.
func (c *ConfigFile[T]) LoadOrCreate(defaultData T) (T, error) {
	data, err := c.Load()
	if err == nil {
		return data, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return data, err
	}

	if err := c.Create(defaultData); err == nil {
		return defaultData, nil
	} else if !errors.Is(err, os.ErrExist) {
		var zero T
		return zero, err
	}
	return c.Load()
}
