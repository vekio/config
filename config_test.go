package config

import (
	"os"
	"path/filepath"
	"testing"
)

type TestConfig struct {
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"database"`
}

// TestConfigPaths verifies the correct construction of application file paths.
func TestConfigPaths(t *testing.T) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("Unable to get user config directory: %v", err)
	}

	// Initialize configuration and check if it reflects the expected values
	conf := NewConfig[TestConfig]()
	expectedAppName := filepath.Base(os.Args[0])
	if conf.AppName() != expectedAppName {
		t.Errorf("Expected AppName to be %s, got %s", expectedAppName, conf.AppName())
	}

	expectedDirPath := filepath.Join(userConfigDir, expectedAppName)
	if conf.DirPath() != expectedDirPath {
		t.Errorf("Expected DirPath to be %s, got %s", expectedDirPath, conf.DirPath())
	}

	expectedFilePath := filepath.Join(expectedDirPath, DefaultConfigFileName)
	if conf.Path() != expectedFilePath {
		t.Errorf("Expected Path to be %s, got %s", expectedFilePath, conf.Path())
	}
}

// TestConfigContentAndInit verifies the initialization and content retrieval from the configuration file.
func TestConfigContentAndInit(t *testing.T) {
	conf := NewConfig[TestConfig]()
	defer os.Remove(conf.Path()) // Ensure cleanup after the test

	// Prepare test data for the configuration initialization
	initialData := TestConfig{}
	initialData.Database.Host = "localhost"
	initialData.Database.Port = 5432

	err := conf.Init(initialData)
	if err != nil {
		t.Fatalf("Failed to initialize configuration file: %v", err)
	}

	// Read and verify the configuration content
	readAndVerifyData(t, conf, initialData)
}

// TestConfigSoftInit verifies the soft initialization of the configuration,
// ensuring it can handle both the creation and loading of configuration data.
func TestConfigSoftInit(t *testing.T) {
	conf := NewConfig[TestConfig]() // Create a new configuration instance for MockConfig
	defer os.Remove(conf.Path())    // Ensure cleanup after the test

	// Prepare test data for the configuration initialization
	initialData := TestConfig{}
	initialData.Database.Host = "localhost"
	initialData.Database.Port = 5432

	// Attempt to load or initialize
	err := conf.Init(initialData)
	if err != nil {
		t.Fatalf("Failed to initialize configuration file: %v", err)
	}

	// Read and verify the configuration data
	readAndVerifyData(t, conf, initialData)

	err = conf.SoftInit()
	if err != nil {
		t.Fatalf("SoftInit failed: %v", err)
	}

	// Read and verify the configuration data
	readAndVerifyData(t, conf, initialData)
}

func readAndVerifyData(t *testing.T, conf *Config[TestConfig], initialData TestConfig) {
	var readData = conf.Data()

	if readData.Database.Host != initialData.Database.Host || readData.Database.Port != initialData.Database.Port {
		t.Errorf("Configuration content does not match expected values")
	}
}
