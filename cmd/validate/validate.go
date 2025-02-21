package validate

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

func NewCmdValidate[T c.Validatable](config *c.ConfigFile[T]) *cli.Command {
	return &cli.Command{
		Name:  "validate",
		Usage: "validate the configuration file",
		Action: func(context.Context, *cli.Command) error {
			data := config.Data()
			if err := data.Validate(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			fmt.Println("configuration is valid.")
			return nil
		},
	}
}
