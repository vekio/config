package config

import (
	"context"
	"fmt"
	"io"
	"strings"

	urfavecli "github.com/urfave/cli/v3"
)

// NewConfigFlag creates a global --config flag that overrides file's path
// when explicitly set by the client application.
func NewConfigFlag[T Validatable](file *ConfigFile[T]) *urfavecli.StringFlag {
	return &urfavecli.StringFlag{
		Name:        "config",
		Usage:       "Path to the configuration file",
		Value:       file.Path(),
		Sources:     urfavecli.EnvVars(configEnvName(file.appName)),
		TakesFile:   true,
		OnlyOnce:    true,
		Config:      urfavecli.StringConfig{TrimSpace: true},
		Validator:   validateConfigFlag,
		Destination: &file.pathOverride,
	}
}

func configEnvName(appName string) string {
	appName = strings.NewReplacer("-", "_", ".", "_").Replace(appName)
	return strings.ToUpper(appName) + "_CONFIG_FILE"
}

func validateConfigFlag(path string) error {
	if _, err := normalizeFilePath(path); err != nil {
		return fmt.Errorf("invalid configuration file path: %w", err)
	}
	return nil
}

// NewConfigCommand creates a reusable config command with show, path, validate,
// and init subcommands. Invoking config without a subcommand displays help.
func NewConfigCommand[T Validatable](file *ConfigFile[T], defaultData T) *urfavecli.Command {
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

func newPathCommand[T Validatable](file *ConfigFile[T]) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "path",
		Usage: "Show the configuration file path",
		Action: func(_ context.Context, cmd *urfavecli.Command) error {
			_, err := fmt.Fprintln(cmd.Root().Writer, file.Path())
			return err
		},
	}
}

func newValidateCommand[T Validatable](file *ConfigFile[T]) *urfavecli.Command {
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

func newInitCommand[T Validatable](file *ConfigFile[T], defaultData T) *urfavecli.Command {
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

func newShowCommand[T Validatable](file *ConfigFile[T]) *urfavecli.Command {
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
