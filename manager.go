package config

// FileManager abstracts how configuration data is loaded from and persisted to
// disk for a given encoding format.
type FileManager[T Validatable] interface {
	// LoadDataFromFile hydrates the provided struct pointer with the contents
	// read from filePath.
	LoadDataFromFile(filePath string, data *T) error
	// WriteDataToFile serializes the given value and stores it at filePath.
	WriteDataToFile(filePath string, data T) error
	// Extension returns the preferred file extension (including the leading
	// dot) used for files handled by this manager.
	Extension() string
}
