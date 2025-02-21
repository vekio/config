package edit

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"

	_file "github.com/vekio/fs/file"
)

func NewCmdEdit[T c.Validatable](config *c.ConfigFile[T]) *cli.Command {
	cmd := &cli.Command{
		Name:  "edit",
		Usage: "edit configuration file",
		Action: func(context.Context, *cli.Command) error {
			err := _file.EditFile(config.Path())
			if err != nil {
				return err
			}
			fmt.Println("configuration file edited successfully.")
			return nil
		},
	}
	return cmd
}
