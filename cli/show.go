package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

// newCmdShow builds the subcommand that prints the raw configuration
// file to stdout so users can quickly inspect the stored values.
func newCmdShow[T c.Validatable]() *cli.Command {
	cmd := &cli.Command{
		Name:        "show",
		Usage:       "Display the current configuration file contents.",
		UsageText:   "conf show",
		Description: "Reads the configuration file from disk and writes its contents to standard output.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			config := c.MustConfigFile[T]()
			buf, err := config.Content()
			if err != nil {
				return fmt.Errorf("read configuration: %w", err)
			}
			if len(buf) == 0 {
				fmt.Fprintln(cmd.Writer, "(configuration file is empty)")
				return nil
			}
			fmt.Fprint(cmd.Writer, string(buf))
			return nil
		},
	}
	return cmd
}
