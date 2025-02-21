package config

type ConfigFileOption[T Validatable] func(*ConfigFile[T])

func WithDefault[T Validatable](defaultData T) ConfigFileOption[T] {
	return func(c *ConfigFile[T]) {
		c.defaultData = defaultData
	}
}
