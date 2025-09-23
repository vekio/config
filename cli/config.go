package cli

import (
	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

// NewCmdConfig wires the configuration management namespace with the
// subcommands that operate on the application's config file.
func NewCmdConfig[T c.Validatable]() *cli.Command {
	cmd := &cli.Command{
		Name:        "conf",
		Usage:       "Manage application's configuration file.",
		UsageText:   "conf [command]",
		Description: "Provides helper commands to show, edit, and validate the configuration file managed by this application.",
		Commands: []*cli.Command{
			newCmdShow[T](),
			newCmdEdit[T](),
			newCmdValidate[T](),
		},
	}
	return cmd
}
