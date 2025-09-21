package config

// ConfigFileOption customizes the behavior of a ConfigFile during construction.
type ConfigFileOption[T Validatable] func(*ConfigFile[T])

// WithDefault seeds the ConfigFile with a default value that will be written to
// disk during initialization when no configuration file exists yet.
func WithDefault[T Validatable](defaultData T) ConfigFileOption[T] {
	return func(c *ConfigFile[T]) {
		if c == nil {
			return
		}
		c.defaultData = defaultData
	}
}

// WithFilename overrides the default file name while ensuring it uses the
// extension that matches the underlying file manager. Environment-specific
// suffixes are re-applied so the behavior mirrors the default constructors.
func WithFilename[T Validatable](fileName string) ConfigFileOption[T] {
	return func(c *ConfigFile[T]) {
		if c == nil {
			return
		}

		extension := ""
		if c.fileManager != nil {
			extension = c.fileManager.Extension()
		}

		normalized := normalizeFileNameWithExtension(fileName, extension)
		if normalized == "" {
			return
		}

		if c.appName != "" {
			normalized = getFileNameForEnvironment(c.appName, normalized)
		}

		c.fileName = normalized
	}
}
