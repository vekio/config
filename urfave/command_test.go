package urfave_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	urfavecli "github.com/urfave/cli/v3"
	config "github.com/vekio/config"
	"github.com/vekio/config/urfave"
)

type testConfig struct {
	Name string `json:"name" yaml:"name"`
	Port int    `json:"port" yaml:"port"`
}

func (c testConfig) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("invalid port")
	}
	return nil
}

func commandDefault() testConfig {
	return testConfig{Name: "default", Port: 8080}
}

func TestConfigShow(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Save(testConfig{Name: "example", Port: 8080}); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	app := &urfavecli.Command{
		Name:     "app",
		Writer:   &output,
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
	}
	if err := app.Run(context.Background(), []string{"app", "config", "show"}); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "name: example\nport: 8080\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestConfigDefaultsToHelp(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	app := &urfavecli.Command{
		Name:     "app",
		Writer:   &output,
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
	}
	if err := app.Run(context.Background(), []string{"app", "config"}); err != nil {
		t.Fatal(err)
	}
	got := output.String()
	for _, name := range []string{"show", "path", "validate", "init"} {
		if !strings.Contains(got, name) {
			t.Fatalf("help output does not list %s:\n%s", name, got)
		}
	}
}

func TestConfigShowMissingFile(t *testing.T) {
	baseDir := t.TempDir()
	file, err := config.NewYAMLConfigFile[testConfig](baseDir, "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}

	app := &urfavecli.Command{
		Name:     "app",
		Writer:   &bytes.Buffer{},
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
	}
	err = app.Run(context.Background(), []string{"app", "config", "show"})
	if err == nil || !strings.Contains(err.Error(), filepath.Join(baseDir, "example", "config.yml")) {
		t.Fatalf("error = %v", err)
	}
}

func TestConfigFlagOverridesPathForSubcommands(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	alternative, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "other", "custom.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err := alternative.Save(testConfig{Name: "alternative", Port: 9090}); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	app := &urfavecli.Command{
		Name:     "app",
		Writer:   &output,
		Flags:    []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
	}
	args := []string{"app", "config", "show", "--config", alternative.Path()}
	if err := app.Run(context.Background(), args); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "name: alternative\nport: 9090\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
	if file.Path() != alternative.Path() {
		t.Fatalf("Path() = %q, want %q", file.Path(), alternative.Path())
	}
}

func TestConfigPath(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	app := &urfavecli.Command{
		Name:     "app",
		Writer:   &output,
		Flags:    []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
	}
	want := filepath.Join(t.TempDir(), "custom.yml")
	if err := app.Run(context.Background(), []string{"app", "config", "path", "--config", want}); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != want+"\n" {
		t.Fatalf("output = %q, want %q", got, want+"\n")
	}
}

func TestConfigValidate(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Save(testConfig{Name: "valid", Port: 8080}); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	app := &urfavecli.Command{
		Name:     "app",
		Writer:   &output,
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
	}
	if err := app.Run(context.Background(), []string{"app", "config", "validate"}); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "configuration is valid\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestConfigInitCreatesExclusively(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	newApp := func() *urfavecli.Command {
		return &urfavecli.Command{
			Name:     "app",
			Writer:   &output,
			Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
		}
	}
	if err := newApp().Run(context.Background(), []string{"app", "config", "init"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := file.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded != commandDefault() {
		t.Fatalf("Load() = %#v, want %#v", loaded, commandDefault())
	}
	if err := newApp().Run(context.Background(), []string{"app", "config", "init"}); !errors.Is(err, os.ErrExist) {
		t.Fatalf("second init error = %v, want os.ErrExist", err)
	}
}

func TestConfigInitForceReplacesExistingFile(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Save(testConfig{Name: "old", Port: 9090}); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	app := &urfavecli.Command{
		Name:     "app",
		Writer:   &output,
		Commands: []*urfavecli.Command{urfave.NewConfigCommand(file, commandDefault())},
	}
	if err := app.Run(context.Background(), []string{"app", "config", "init", "--force"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := file.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded != commandDefault() {
		t.Fatalf("Load() = %#v, want %#v", loaded, commandDefault())
	}
}

func TestConfigFlagRejectsEmptyPath(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	app := &urfavecli.Command{
		Name:  "app",
		Flags: []urfavecli.Flag{urfave.NewConfigFlag(file)},
	}
	if err := app.Run(context.Background(), []string{"app", "--config", " "}); err == nil {
		t.Fatal("--config accepted an empty path")
	}
}

func TestConfigFlagDocumentsEnvironmentVariable(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "my-app", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	app := &urfavecli.Command{
		Name:   "app",
		Writer: &output,
		Flags:  []urfavecli.Flag{urfave.NewConfigFlag(file)},
	}
	if err := app.Run(context.Background(), []string{"app", "--help"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "MY_APP_CONFIG_FILE") {
		t.Fatalf("help output does not document MY_APP_CONFIG_FILE:\n%s", output.String())
	}
}

func TestConfigFlagDefersFilesystemValidation(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	app := &urfavecli.Command{
		Name:  "app",
		Flags: []urfavecli.Flag{urfave.NewConfigFlag(file)},
	}
	if err := app.Run(context.Background(), []string{"app", "--config", directory}); err != nil {
		t.Fatalf("--config performed filesystem validation: %v", err)
	}
	if file.Path() != directory {
		t.Fatalf("Path() = %q, want %q", file.Path(), directory)
	}
}

func TestConfigFlagAcceptsFileThatDoesNotExistYet(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(t.TempDir(), "new", "config.yml")
	app := &urfavecli.Command{
		Name:   "app",
		Flags:  []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Action: func(context.Context, *urfavecli.Command) error { return nil },
	}
	if err := app.Run(context.Background(), []string{"app", "--config", want}); err != nil {
		t.Fatal(err)
	}
	if file.Path() != want {
		t.Fatalf("Path() = %q, want %q", file.Path(), want)
	}
}

func TestConfigFlagOverridesPathBeforeBeforeHook(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	alternative, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "other", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	want := testConfig{Name: "alternative", Port: 9090}
	if err := alternative.Save(want); err != nil {
		t.Fatal(err)
	}

	var loaded testConfig
	app := &urfavecli.Command{
		Name:  "app",
		Flags: []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Before: func(ctx context.Context, _ *urfavecli.Command) (context.Context, error) {
			loaded, err = file.Load()
			return ctx, err
		},
		Commands: []*urfavecli.Command{
			{
				Name:   "serve",
				Action: func(context.Context, *urfavecli.Command) error { return nil },
			},
		},
	}
	if err := app.Run(context.Background(), []string{"app", "--config", alternative.Path(), "serve"}); err != nil {
		t.Fatal(err)
	}
	if loaded != want {
		t.Fatalf("Before loaded %#v, want %#v", loaded, want)
	}
}

func TestConfigFlagUsesEnvironmentSource(t *testing.T) {
	want := filepath.Join(t.TempDir(), "environment.yml")
	t.Setenv("MY_APP_CONFIG_FILE", want)
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "my-app", "config.yml")
	if err != nil {
		t.Fatal(err)
	}

	app := &urfavecli.Command{
		Name:   "app",
		Flags:  []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Action: func(context.Context, *urfavecli.Command) error { return nil },
	}
	if err := app.Run(context.Background(), []string{"app"}); err != nil {
		t.Fatal(err)
	}
	if file.Path() != want {
		t.Fatalf("Path() = %q, want %q", file.Path(), want)
	}
}

func TestConfigFlagTakesPrecedenceOverEnvironment(t *testing.T) {
	environmentPath := filepath.Join(t.TempDir(), "environment.yml")
	flagPath := filepath.Join(t.TempDir(), "flag.yml")
	t.Setenv("MY_APP_CONFIG_FILE", environmentPath)
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "my.app", "config.yml")
	if err != nil {
		t.Fatal(err)
	}

	app := &urfavecli.Command{
		Name:   "app",
		Flags:  []urfavecli.Flag{urfave.NewConfigFlag(file)},
		Action: func(context.Context, *urfavecli.Command) error { return nil },
	}
	if err := app.Run(context.Background(), []string{"app", "--config", flagPath}); err != nil {
		t.Fatal(err)
	}
	if file.Path() != flagPath {
		t.Fatalf("Path() = %q, want %q", file.Path(), flagPath)
	}
}
