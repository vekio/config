package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vekio/config"
)

func main() {
	config, err := NewDatabaseConfig()
	if err != nil {
		log.Fatalf("error creating database config %v", err)
	}
	defer os.RemoveAll(config.DirPath())

	// fmt.Println(config.Host())

	buf, err := config.Content()
	if err != nil {
		log.Fatalf("error reading content: %v", err)
	}
	fmt.Println(string(buf))
}

type Database struct {
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"database"`
}

func (d Database) Validate() error {
	return nil
}

type DatabaseConfig struct {
	config.ConfigFile[Database]
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	database := Database{}
	database.Database.Host = "localhost"
	database.Database.Port = 987

	config, err := config.NewDefaultConfigFile[Database]("database", config.WithDefault(database))
	if err != nil {
		return nil, fmt.Errorf("error creating DefaultConfigFile: %w", err)
	}

	err = config.SoftInit()
	if err != nil {
		return nil, fmt.Errorf("error initializing or loading config: %w", err)
	}

	c := &DatabaseConfig{
		ConfigFile: *config,
	}
	return c, nil
}

func (c *DatabaseConfig) Host() string {
	return c.Data().Database.Host
}
