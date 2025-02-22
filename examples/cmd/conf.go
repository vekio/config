package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vekio/config"
	cmdConf "github.com/vekio/config/cmd"
)

func main() {
	config, err := NewDatabaseConfig()
	if err != nil {
		log.Fatalf("error creating database config %v", err)
	}
	defer os.RemoveAll(config.DirPath())

	cmd := cmdConf.NewCmdConfig(config.ConfigFile)
	cmd.Run(context.Background(), os.Args)
}

type Settings struct {
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"database"`
}

func (s Settings) Validate() error {
	if s.Database.Port != 5432 {
		return fmt.Errorf("database port unknown: %d", s.Database.Port)
	}
	return nil
}

type DatabaseConfig struct {
	config.ConfigFile[Settings]
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	database := Settings{}
	database.Database.Host = "localhost"
	database.Database.Port = 987

	config, err := config.NewDefaultConfigFile[Settings]("database", config.WithDefault(database))
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
