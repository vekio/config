package urfave_test

import (
	"fmt"

	urfavecli "github.com/urfave/cli/v3"
	"github.com/vekio/config"
	"github.com/vekio/config/urfave"
)

type exampleConfig struct{}

func (exampleConfig) Validate() error { return nil }

func Example() {
	file, err := config.NewYAMLConfigFile[exampleConfig]("example", "config.yml")
	if err != nil {
		panic(err)
	}

	app := &urfavecli.Command{
		Name:     "example",
		Flags:    []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file)},
	}

	fmt.Println(app.Name, app.Flags[0].Names()[0], app.Commands[0].Name)
	// Output: example config config
}
