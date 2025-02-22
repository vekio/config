package cmd

import (
	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

func NewCmdConfig[T c.Validatable](config c.ConfigFile[T]) *cli.Command {
	cmd := &cli.Command{
		Name:  "conf",
		Usage: "manage configuration",
		Commands: []*cli.Command{
			NewCmdShow(config),
			NewCmdEdit(config),
			NewCmdValidate(config),
		},
	}
	return cmd
}
