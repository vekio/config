package urfave

import (
	"context"
	"fmt"
	"io"

	urfavecli "github.com/urfave/cli/v3"
	"github.com/vekio/config"
)

// NewConfigFlag creates a global --config flag that overrides file's path
// when explicitly set by the client application.
func NewConfigFlag[T config.Validatable](file *config.ConfigFile[T]) urfavecli.Flag {
	flag := &urfavecli.StringFlag{
		Name:        "config",
		Usage:       "Path to the configuration file",
		Value:       file.Path(),
		TakesFile:   true,
		OnlyOnce:    true,
		HideDefault: true,
		Config:      urfavecli.StringConfig{TrimSpace: true},
	}
	return &configFlag[T]{StringFlag: flag, file: file}
}

// configFlag synchronizes parsed values before command Before hooks run.
// Embedding preserves urfave's help and flag metadata interfaces.
type configFlag[T config.Validatable] struct {
	*urfavecli.StringFlag
	file *config.ConfigFile[T]
}

// Set synchronizes command-line values, including global flags placed after a
// subcommand, as soon as urfave parses them. Environment sources are applied
// internally by StringFlag and have already been resolved by ConfigFile, so
// they do not dispatch through this method or update the path twice.
func (f *configFlag[T]) Set(name, value string) error {
	if err := f.StringFlag.Set(name, value); err != nil {
		return err
	}
	return f.file.SetPath(value)
}

// NewConfigCommand creates a reusable config command with show, path, validate,
// and init subcommands. Invoking config without a subcommand displays help.
func NewConfigCommand[T config.Validatable](file *config.ConfigFile[T], defaultData T) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "config",
		Usage: "Manage the application configuration",
		Commands: []*urfavecli.Command{
			newShowCommand(file),
			newPathCommand(file),
			newValidateCommand(file),
			newInitCommand(file, defaultData),
		},
	}
}

func newPathCommand[T config.Validatable](file *config.ConfigFile[T]) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "path",
		Usage: "Show the configuration file path",
		Action: func(_ context.Context, cmd *urfavecli.Command) error {
			_, err := fmt.Fprintln(cmd.Root().Writer, file.Path())
			return err
		},
	}
}

func newValidateCommand[T config.Validatable](file *config.ConfigFile[T]) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "validate",
		Usage: "Validate the configuration file",
		Action: func(_ context.Context, cmd *urfavecli.Command) error {
			if _, err := file.Load(); err != nil {
				return fmt.Errorf("validate configuration: %w", err)
			}
			_, err := fmt.Fprintln(cmd.Root().Writer, "configuration is valid")
			return err
		},
	}
}

func newInitCommand[T config.Validatable](file *config.ConfigFile[T], defaultData T) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "init",
		Usage: "Create the default configuration file",
		Flags: []urfavecli.Flag{
			&urfavecli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Replace the configuration file if it already exists",
			},
		},
		Action: func(_ context.Context, cmd *urfavecli.Command) error {
			write := file.Create
			if cmd.Bool("force") {
				write = file.Save
			}
			if err := write(defaultData); err != nil {
				return fmt.Errorf("initialize configuration: %w", err)
			}
			_, err := fmt.Fprintln(cmd.Root().Writer, file.Path())
			return err
		},
	}
}

func newShowCommand[T config.Validatable](file *config.ConfigFile[T]) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "show",
		Usage: "Show the configuration file contents",
		Action: func(_ context.Context, cmd *urfavecli.Command) error {
			content, err := file.Content()
			if err != nil {
				return fmt.Errorf("show configuration: %w", err)
			}
			if _, err := cmd.Root().Writer.Write(content); err != nil {
				return fmt.Errorf("write configuration: %w", err)
			}
			if len(content) == 0 || content[len(content)-1] != '\n' {
				if _, err := io.WriteString(cmd.Root().Writer, "\n"); err != nil {
					return fmt.Errorf("write configuration: %w", err)
				}
			}
			return nil
		},
	}
}
