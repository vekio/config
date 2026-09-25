package config

// Option configures a ConfigFile during construction.
type Option[T Validatable] func(*ConfigFile[T])

// Default sets the value used by ConfigFile.LoadOrCreate when the
// configuration file does not exist. Without this option, the zero value of T
// is used.
func Default[T Validatable](value T) Option[T] {
	return func(file *ConfigFile[T]) {
		file.defaults = value
	}
}
