package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

// newCmdValidate provides the subcommand that reloads the configuration and
// executes the user-defined validation logic.
func newCmdValidate[T c.Validatable]() *cli.Command {
	return &cli.Command{
		Name:        "validate",
		Usage:       "Reload and validate the configuration contents.",
		UsageText:   "conf validate",
		Description: "Reloads the configuration file from disk and runs the Validate method implemented by the consumer-provided config struct.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			config := c.MustConfigFile[T]()
			if err := config.Data().Validate(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			fmt.Fprintln(cmd.Writer, "configuration is valid")
			return nil
		},
	}
}
