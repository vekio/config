package config

import (
	"encoding/json"
	"fmt"

	_file "github.com/vekio/fs/file"
)

// JSONFileManager implements IBuilder for JSON files
type JSONFileManager[T any] struct{}

func NewJSONFileManager[T any]() *JSONFileManager[T] {
	return &JSONFileManager[T]{}
}

func (b *JSONFileManager[T]) LoadDataFromFile(filePath string, data *T) error {
	buf, err := _file.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buf, data); err != nil {
		return fmt.Errorf("error unmarshaling JSON data: %w", err)
	}
	return nil
}

func (b *JSONFileManager[T]) WriteDataToFile(filePath string, data T) error {
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON data: %w", err)
	}
	if err := _file.WriteFileContent(filePath, buf, _file.RestrictedFilePerms); err != nil {
		return fmt.Errorf("error writing JSON data to file: %w", err)
	}
	return nil
}
