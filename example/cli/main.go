package main

import (
	"context"
	"fmt"
	"log"
	"os"

	urfavecli "github.com/urfave/cli/v3"
	"github.com/vekio/config"
)

type Config struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func (c Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func main() {
	configFile, err := config.NewDefaultConfigFile[Config]("config-example")
	if err != nil {
		log.Fatal(err)
	}

	var current Config
	defaultConfig := Config{
		Address: "127.0.0.1",
		Port:    8080,
	}
	app := &urfavecli.Command{
		Name:  "example",
		Usage: "Example application with configurable file path",
		Flags: []urfavecli.Flag{config.NewConfigFlag(configFile)},
		Commands: []*urfavecli.Command{
			config.NewConfigCommand(configFile, defaultConfig),
			{
				Name:  "serve",
				Usage: "Show the address that the application would use",
				Before: func(ctx context.Context, _ *urfavecli.Command) (context.Context, error) {
					loaded, err := configFile.LoadOrCreate(defaultConfig)
					if err != nil {
						return ctx, err
					}
					current = loaded
					return ctx, nil
				},
				Action: func(_ context.Context, _ *urfavecli.Command) error {
					fmt.Printf("Serving on %s:%d\n", current.Address, current.Port)
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
