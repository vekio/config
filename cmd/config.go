package cmd

import (
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
	cmdEdit "github.com/vekio/config/cmd/edit"
	cmdShow "github.com/vekio/config/cmd/show"
)

func NewCmdConfig(config *c.Config) *cli.Command {
	cmd := &cli.Command{
		Name:  "conf",
		Usage: fmt.Sprintf("configuration for %s", config.AppName()),
		Commands: []*cli.Command{
			cmdShow.NewCmdShow(config),
			cmdEdit.NewCmdEdit(config),
		},
	}

	return cmd
}
