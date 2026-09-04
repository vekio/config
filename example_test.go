package config_test

import (
	"fmt"

	"github.com/vekio/config"
)

type exampleConfig struct {
	Address string `yaml:"address"`
}

func (c exampleConfig) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}
	return nil
}

func ExampleNewYAMLConfigFile() {
	file, err := config.NewYAMLConfigFile[exampleConfig]("example", "config.yml")
	if err != nil {
		panic(err)
	}
	if err := file.SetPath("./config.dev.yml"); err != nil {
		panic(err)
	}

	fmt.Println(file.Path())
	// Output: config.dev.yml
}
