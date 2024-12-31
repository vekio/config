package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

type MockConfig struct {
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
	conf := NewConfig()
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
	conf := NewConfig()
	defer os.Remove(conf.Path()) // Ensure cleanup after the test

	// Prepare test data for the configuration initialization
	initialData := MockConfig{}
	initialData.Database.Host = "localhost"
	initialData.Database.Port = 5432

	err := conf.Init(&initialData)
	if err != nil {
		t.Fatalf("Failed to initialize configuration file: %v", err)
	}

	// Read and verify the configuration content
	readAndVerifiyContent(t, conf, initialData)
}

// TestConfigSoftInit verifies the soft initialization of the configuration,
// ensuring it can handle both the creation and loading of configuration data.
func TestConfigSoftInit(t *testing.T) {
	conf := NewConfig()
	defer os.Remove(conf.Path()) // Ensure cleanup after the test

	// Ensure the configuration file does not exist initially
	os.Remove(conf.Path())

	// Initial data setup for the test
	initialData := MockConfig{}
	initialData.Database.Host = "localhost"
	initialData.Database.Port = 5432

	err := conf.SoftInit(&initialData)
	if err != nil {
		t.Fatalf("Failed to soft initialize configuration file: %v", err)
	}

	// Read and verify the configuration content
	readAndVerifiyContent(t, conf, initialData)

	// Setup for verifying soft data loading behavior
	modifiedData := MockConfig{}
	modifiedData.Database.Host = "localhost"
	modifiedData.Database.Port = 3306 // Change of port to verify loading

	err = conf.SoftInit(&modifiedData)
	if err != nil {
		t.Fatalf("Failed to soft initialize existing configuration file: %v", err)
	}

	// Ensure that the data is consistent and was not altered unexpectedly
	if modifiedData.Database.Port != initialData.Database.Port {
		t.Errorf("Expected port to be %v, got %v", 5432, modifiedData.Database.Port)
	}
}

func readAndVerifiyContent(t *testing.T, conf *Config, initialData MockConfig) {
	content, err := conf.Content()
	if err != nil {
		t.Fatalf("Failed to read configuration content: %v", err)
	}

	var readData MockConfig
	err = yaml.Unmarshal(content, &readData)
	if err != nil {
		t.Fatalf("Failed to unmarshal configuration content: %v", err)
	}

	if readData.Database.Host != initialData.Database.Host || readData.Database.Port != initialData.Database.Port {
		t.Errorf("Configuration content does not match expected values")
	}
}
