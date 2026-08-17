package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func defaultConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user configuration directory: %w", err)
	}
	configDir, err = normalizeDirPath(configDir)
	if err != nil {
		return "", fmt.Errorf("resolve user configuration directory: %w", err)
	}
	return configDir, nil
}

// DefaultDataDir returns the application's conventional data directory.
// XDG_DATA_HOME takes precedence when set and must contain an absolute path.
func DefaultDataDir(appName string) (string, error) {
	if !isPathElement(appName) {
		return "", fmt.Errorf("application name %q must be a single path element", appName)
	}

	var baseDir string
	if dataDir := os.Getenv("XDG_DATA_HOME"); dataDir != "" {
		if !filepath.IsAbs(dataDir) {
			return "", fmt.Errorf("XDG_DATA_HOME must be an absolute path")
		}
		baseDir = dataDir
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".local", "share")
	}
	baseDir, err := normalizeDirPath(baseDir)
	if err != nil {
		return "", fmt.Errorf("resolve user data directory: %w", err)
	}
	return filepath.Join(baseDir, appName), nil
}
