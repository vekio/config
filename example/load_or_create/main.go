package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/vekio/config"
)

type Config struct {
	Address  string `json:"address" yaml:"address"`
	Port     int    `json:"port" yaml:"port"`
	DataPath string `json:"data_path" yaml:"data_path"`
}

func (c Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.DataPath == "" {
		return fmt.Errorf("data path is required")
	}
	return nil
}

func main() {
	dataDir, err := config.DefaultDataDir("config-example")
	if err != nil {
		log.Fatal(err)
	}
	defaults := Config{
		Address:  "127.0.0.1",
		Port:     8080,
		DataPath: filepath.Join(dataDir, "data.db"),
	}
	configFile, err := config.NewYAMLConfigFile(
		"config-example",
		"config.yml",
		config.Default(defaults),
	)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := configFile.LoadOrCreate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Configuration file: %s\n", configFile.Path())
	fmt.Printf("Configuration: %+v\n", cfg)
}
