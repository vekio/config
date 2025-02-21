package cmd

import (
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
	cmdEdit "github.com/vekio/config/cmd/edit"
	cmdShow "github.com/vekio/config/cmd/show"
	cmdValidate "github.com/vekio/config/cmd/validate"
)

func NewCmdConfig[T c.Validatable](config *c.ConfigFile[T]) *cli.Command {
	cmd := &cli.Command{
		Name:  "conf",
		Usage: fmt.Sprintf("manage configuration"),
		Commands: []*cli.Command{
			cmdShow.NewCmdShow(config),
			cmdEdit.NewCmdEdit(config),
			cmdValidate.NewCmdValidate(config),
		},
	}
	return cmd
}
