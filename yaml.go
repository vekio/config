package config

import (
	"fmt"

	_file "github.com/vekio/fs/file"
	"gopkg.in/yaml.v3"
)

// YAMLFileManager implements IBuilder for YAML files
type YAMLFileManager[T any] struct{}

func NewYAMLFileManager[T any]() *YAMLFileManager[T] {
	return &YAMLFileManager[T]{}
}

func (b *YAMLFileManager[T]) LoadDataFromFile(filePath string, data T) error {
	buf, err := _file.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(buf, &data); err != nil {
		return fmt.Errorf("error unmarshaling YAML data: %w", err)
	}
	return nil
}

func (b *YAMLFileManager[T]) WriteDataToFile(filePath string, data T) error {
	buf, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling YAML data: %w", err)
	}
	if err := _file.WriteFileContent(filePath, buf, _file.RestrictedFilePerms); err != nil {
		return fmt.Errorf("error writing YAML data to file: %w", err)
	}
	return nil
}
