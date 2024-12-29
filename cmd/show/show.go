package show

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

func NewCmdShow(config *c.Config) *cli.Command {
	cmd := &cli.Command{
		Name:  "edit",
		Usage: "edit configuration file",
		Action: func(context.Context, *cli.Command) error {
			buf, err := config.Content()
			if err != nil {
				return err
			}
			fmt.Print(string(buf))
			return nil
		},
	}

	return cmd
}
