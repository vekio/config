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

// NewCLIConfig builds a CLI configuration helper for the provided application.
func NewCLIConfig[T c.Validatable](appName string, defaults T, options ...c.ConfigFileOption[T]) (*CLIConfig[T], error) {
	opts := make([]c.ConfigFileOption[T], 0, len(options)+1)
	opts = append(opts, c.WithDefault(defaults))
	opts = append(opts, options...)

	file, err := c.NewDefaultConfigFile(appName, opts...)
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
