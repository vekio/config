package config

import (
	"fmt"
	"os"

	"github.com/vekio/x/fs"
	"gopkg.in/yaml.v3"
)

// YAMLFileManager implements FileManager for configurations encoded as YAML.
type YAMLFileManager[T any] struct{}

// NewYAMLFileManager builds a FileManager that marshals and unmarshals YAML
// payloads via gopkg.in/yaml.v3.
func NewYAMLFileManager[T any]() *YAMLFileManager[T] {
	return &YAMLFileManager[T]{}
}

// Extension returns the canonical YAML file extension.
func (b *YAMLFileManager[T]) Extension() string {
	return ".yml"
}

// LoadDataFromFile reads the YAML file, unmarshals it into the provided value,
// and returns an error if the file cannot be read or parsed.
func (b *YAMLFileManager[T]) LoadDataFromFile(filePath string, data *T) error {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read YAML file: %w", err)
	}
	if err := yaml.Unmarshal(buf, data); err != nil {
		return fmt.Errorf("error unmarshaling YAML data: %w", err)
	}
	return nil
}

// WriteDataToFile serializes the value as YAML and persists it to disk.
func (b *YAMLFileManager[T]) WriteDataToFile(filePath string, data T) error {
	buf, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling YAML data: %w", err)
	}
	// Ensure the destination exists and persists the encoded payload.
	if err := fs.WriteFileWithDirs(filePath, buf, fs.RestrictedFileMode); err != nil {
		return fmt.Errorf("error writing YAML data to file: %w", err)
	}
	return nil
}
