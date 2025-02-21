package config

type FileManager[T Validatable] interface {
	LoadDataFromFile(filePath string, data T) error
	WriteDataToFile(filePath string, data T) error
}
