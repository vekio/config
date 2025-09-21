package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vekio/x/fs"
)

// JSONFileManager implements FileManager for configurations encoded as JSON.
type JSONFileManager[T any] struct{}

// NewJSONFileManager builds a FileManager that marshals and unmarshals JSON
// payloads using the standard library.
func NewJSONFileManager[T any]() *JSONFileManager[T] {
	return &JSONFileManager[T]{}
}

// Extension returns the canonical JSON file extension.
func (b *JSONFileManager[T]) Extension() string {
	return ".json"
}

// LoadDataFromFile reads the JSON file, unmarshals it into the provided value,
// and returns an error if the file cannot be read or parsed.
func (b *JSONFileManager[T]) LoadDataFromFile(filePath string, data *T) error {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read JSON file: %w", err)
	}
	if err := json.Unmarshal(buf, data); err != nil {
		return fmt.Errorf("error unmarshaling JSON data: %w", err)
	}
	return nil
}

// WriteDataToFile serializes the value as JSON and persists it to disk.
func (b *JSONFileManager[T]) WriteDataToFile(filePath string, data T) error {
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON data: %w", err)
	}
	// Ensure the destination exists and persists the encoded payload.
	if err := fs.WriteFileWithDirs(filePath, buf, fs.RestrictedFileMode); err != nil {
		return fmt.Errorf("error writing JSON data to file: %w", err)
	}
	return nil
}
