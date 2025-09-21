package cli

import (
	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

// CLIConfig groups the configuration file and CLI command for configuration management.
type CLIConfig[T c.Validatable] struct {
	File    *c.ConfigFile[T]
	Command *cli.Command
}

// AddSubcommands appends additional CLI commands under the "conf" namespace.
func (c *CLIConfig[T]) AddSubcommands(commands ...*cli.Command) {
	if c == nil || c.Command == nil || len(commands) == 0 {
		return
	}

	for _, command := range commands {
		if command == nil {
			continue
		}
		c.Command.Commands = append(c.Command.Commands, command)
	}
}

// NewCLIConfig builds a CLI configuration helper for the provided application.
// The application name defaults to the executable name but can be overridden
// with config.WithAppName.
func NewCLIConfig[T c.Validatable](defaults T, options ...c.ConfigFileOption[T]) (*CLIConfig[T], error) {
	opts := make([]c.ConfigFileOption[T], 0, len(options)+1)
	opts = append(opts, c.WithDefault(defaults))
	opts = append(opts, options...)

	file, err := c.NewDefaultConfigFile(opts...)
	if err != nil {
		return nil, err
	}

	if err := file.SoftInit(); err != nil {
		return nil, err
	}

	cmd := newCmdConfig(file)

	return &CLIConfig[T]{
		File:    file,
		Command: cmd,
	}, nil
}

// Attach registers the configuration command into the provided CLI application.
func (c *CLIConfig[T]) Attach(app *cli.Command) {
	if c == nil || app == nil || c.Command == nil {
		return
	}
	app.Commands = append(app.Commands, c.Command)
}
