package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// cleanPath performs only syntactic path normalization. Filesystem properties
// are checked by the operation that eventually reads or writes the path.
func cleanPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	return filepath.Clean(path), nil
}

func configPathEnvName(appName string) string {
	var name strings.Builder
	for _, char := range appName {
		switch {
		case char >= 'a' && char <= 'z':
			name.WriteRune(char - ('a' - 'A'))
		case char >= 'A' && char <= 'Z', char >= '0' && char <= '9', char == '_':
			name.WriteRune(char)
		default:
			name.WriteByte('_')
		}
	}
	if value := name.String(); value[0] >= '0' && value[0] <= '9' {
		return "_" + value + "_CONFIG_FILE"
	}
	return name.String() + "_CONFIG_FILE"
}
