package config

import (
	"os"
	"path/filepath"
	"strings"
)

func defaultConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	return dir
}

func defaultAppName() string {
	if len(os.Args) == 0 {
		return "app"
	}

	name := strings.TrimSpace(filepath.Base(os.Args[0]))
	if name == "" {
		return "app"
	}

	return name
}
