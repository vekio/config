package edit

import (
	"context"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"

	_file "github.com/vekio/fs/file"
)

func NewCmdEdit[T any](config *c.Config[T]) *cli.Command {
	cmd := &cli.Command{
		Name:  "edit",
		Usage: "edit configuration file",
		Action: func(context.Context, *cli.Command) error {
			err := _file.EditFile(config.Path())
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
