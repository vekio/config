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
