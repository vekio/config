package show

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

func NewCmdShow[T c.Validatable](config *c.ConfigFile[T]) *cli.Command {
	cmd := &cli.Command{
		Name:  "show",
		Usage: "show configuration file content",
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
