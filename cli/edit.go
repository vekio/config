package cli

import (
	"context"

	"github.com/urfave/cli/v3"
	c "github.com/vekio/config"
)

// newCmdEdit registers the subcommand that opens the managed configuration
// file in the user's preferred editor. It ensures the file exists before
// launching the editor so the command can operate on fresh installations.
func newCmdEdit[T c.Validatable]() *cli.Command {
	cmd := &cli.Command{
		Name:        "edit",
		Usage:       "Open the configuration file in the preferred editor.",
		UsageText:   "conf edit",
		Description: "Creates the configuration file if it does not exist and then launches $VISUAL, $EDITOR, nano, or vi in that order.",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c.MustConfigFile[T]()
			// TODO create fs/file Edit()
			// if err := file.Touch(config.Path(), fs.DefaultFileMode); err != nil {
			// 	return fmt.Errorf("prepare configuration file: %w", err)
			// }

			// editor, err := selectEditor()
			// if err != nil {
			// 	return fmt.Errorf("determine editor: %w", err)
			// }

			// if err := launchEditor(ctx, editor, config.Path()); err != nil {
			// 	return err
			// }

			// fmt.Fprintln(cmd.Writer, "configuration file opened in editor")
			return nil
		},
	}
	return cmd
}

// func selectEditor() ([]string, error) {
// 	candidates := make([]string, 0, 4)
// 	for _, key := range []string{"VISUAL", "EDITOR"} {
// 		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
// 			candidates = append(candidates, value)
// 		}
// 	}
// 	candidates = append(candidates, "nano", "vi")

// 	for _, candidate := range candidates {
// 		parts := strings.Fields(candidate)
// 		if len(parts) == 0 {
// 			continue
// 		}

// 		if _, err := exec.LookPath(parts[0]); err != nil {
// 			continue
// 		}

// 		return parts, nil
// 	}

// 	return nil, errors.New("no editor available: set $VISUAL or $EDITOR, or install nano/vi")
// }

// func launchEditor(ctx context.Context, editor []string, path string) error {
// 	if len(editor) == 0 {
// 		return errors.New("editor command is empty")
// 	}

// 	cmd := exec.CommandContext(ctx, editor[0], append(editor[1:], path)...)
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("launch editor: %w", err)
// 	}
// 	return nil
// }
