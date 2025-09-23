package config

import (
	"fmt"
	"strings"
)

// Format identifies the serialization format used to persist the configuration file.
type Format string

const (
	// FormatYAML stores configuration as YAML.
	FormatYAML Format = "yaml"
	// FormatJSON stores configuration as JSON.
	FormatJSON Format = "json"
)

func (f Format) normalized() string {
	return strings.ToLower(strings.TrimSpace(string(f)))
}

// NewConfigFile constructs a ConfigFile using the requested format. When format is empty it defaults to YAML.
func NewConfigFile[T Validatable](format Format, options ...ConfigFileOption[T]) (*ConfigFile[T], error) {
	switch normalized := format.normalized(); normalized {
	case "", string(FormatYAML), "yml":
		return NewYAMLConfigFile(options...)
	case string(FormatJSON):
		return NewJSONConfigFile(options...)
	default:
		return nil, fmt.Errorf("config: unsupported format %q", normalized)
	}
}
