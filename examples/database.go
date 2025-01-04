package main

import (
	"fmt"
	"log"

	"github.com/vekio/config"
)

func main() {
	config, err := NewDatabaseConfig()
	if err != nil {
		log.Fatalf("failed creating database config %v", err)
	}

	fmt.Println(config.Host())

	buf, _ := config.Content()
	fmt.Println(string(buf))
}

type Database struct {
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"database"`
}

type DatabaseConfig struct {
	config.Config[Database]
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	config := config.NewConfig[Database]()

	err := config.SoftInit()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize or load config: %v", err)
	}

	c := &DatabaseConfig{
		Config: *config,
	}
	return c, nil
}

func (c *DatabaseConfig) Host() string {
	return c.Data().Database.Host
}
