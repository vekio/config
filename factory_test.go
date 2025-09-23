package config

import "testing"

func TestNewConfigFileDefaultsToYAML(t *testing.T) {
	cfg, err := NewConfigFile[testSettings]("")
	if err != nil {
		t.Fatalf("NewConfigFile returned error: %v", err)
	}

	if cfg == nil {
		t.Fatalf("expected config file, got nil")
	}

	if cfg.fileManager.Extension() != ".yml" {
		t.Fatalf("expected YAML extension, got %q", cfg.fileManager.Extension())
	}
}

func TestNewConfigFileYAML(t *testing.T) {
	cfg, err := NewConfigFile[testSettings](FormatYAML)
	if err != nil {
		t.Fatalf("NewConfigFile returned error: %v", err)
	}

	if cfg.fileManager.Extension() != ".yml" {
		t.Fatalf("expected YAML extension, got %q", cfg.fileManager.Extension())
	}
}

func TestNewConfigFileJSON(t *testing.T) {
	cfg, err := NewConfigFile[testSettings](FormatJSON)
	if err != nil {
		t.Fatalf("NewConfigFile returned error: %v", err)
	}

	if cfg.fileManager.Extension() != ".json" {
		t.Fatalf("expected JSON extension, got %q", cfg.fileManager.Extension())
	}
}

func TestNewConfigFileInvalidFormat(t *testing.T) {
	if _, err := NewConfigFile[testSettings]("toml"); err == nil {
		t.Fatalf("expected error for unsupported format")
	}
}
