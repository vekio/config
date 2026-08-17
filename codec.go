package config

// codec converts typed configuration values to and from their on-disk
// representation. ConfigFile owns all filesystem operations.
type codec[T any] interface {
	marshal(data T) ([]byte, error)
	unmarshal(content []byte, data *T) error
}
