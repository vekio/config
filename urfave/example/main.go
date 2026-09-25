// A small item application showing configuration, CLI integration, and
// runtime selection of an in-memory or local filesystem repository.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	urfavecli "github.com/urfave/cli/v3"
	"github.com/vekio/config"
	"github.com/vekio/config/urfave"
)

func main() {
	if err := run(context.Background(), os.Args, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	dataDir, err := config.DefaultDataDir("items-example")
	if err != nil {
		return err
	}
	file, err := config.NewYAMLConfigFile(
		"items-example", "config.yml",
		config.Default(Config{Storage: "localfs", DataDir: dataDir}),
	)
	if err != nil {
		return err
	}

	var app *App
	command := &urfavecli.Command{
		Name:   "items",
		Usage:  "Store and list items",
		Writer: out,
		Flags:  []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Commands: []*urfavecli.Command{
			urfave.NewConfigCommand(file),
			{
				Name:  "item",
				Usage: "Work with items",
				Before: func(ctx context.Context, _ *urfavecli.Command) (context.Context, error) {
					settings, err := file.LoadOrCreate()
					if err != nil {
						return ctx, err
					}
					app, err = NewApp(settings)
					return ctx, err
				},
				Commands: []*urfavecli.Command{
					{
						Name:  "add",
						Usage: "Add an item",
						Flags: []urfavecli.Flag{&urfavecli.StringFlag{Name: "name", Required: true}},
						Action: func(_ context.Context, cmd *urfavecli.Command) error {
							item, err := app.AddItem.Run(cmd.String("name"))
							if err != nil {
								return err
							}
							_, err = fmt.Fprintf(out, "%s %s\n", item.ID, item.Name)
							return err
						},
					},
					{
						Name:  "list",
						Usage: "List items",
						Action: func(_ context.Context, _ *urfavecli.Command) error {
							items, err := app.ListItems.Run()
							if err != nil {
								return err
							}
							for _, item := range items {
								if _, err := fmt.Fprintf(out, "%s %s\n", item.ID, item.Name); err != nil {
									return err
								}
							}
							return nil
						},
					},
				},
			},
		},
	}
	return command.Run(ctx, args)
}
