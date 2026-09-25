package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppSelectsRepository(t *testing.T) {
	dataDir := filepath.Join(t.TempDir(), "items")
	for _, storage := range []string{"memory", "localfs"} {
		t.Run(storage, func(t *testing.T) {
			cfg := Config{Storage: storage, DataDir: dataDir}
			app, err := NewApp(cfg)
			if err != nil {
				t.Fatal(err)
			}
			created, err := app.AddItem.Run("first")
			if err != nil {
				t.Fatal(err)
			}
			items, err := app.ListItems.Run()
			if err != nil {
				t.Fatal(err)
			}
			if len(items) != 1 || items[0] != created {
				t.Fatalf("ListItems.Run() = %#v, want %#v", items, created)
			}
			second, err := NewApp(cfg)
			if err != nil {
				t.Fatal(err)
			}
			items, err = second.ListItems.Run()
			if err != nil {
				t.Fatal(err)
			}
			if storage == "memory" && len(items) != 0 {
				t.Fatalf("new memory app retained items: %#v", items)
			}
			if storage == "localfs" && (len(items) != 1 || items[0] != created) {
				t.Fatalf("new localfs app lost item: %#v", items)
			}
		})
	}
}

func TestCLIFlagOverridesEnvironmentBeforeCreatingApp(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_DATA_HOME", filepath.Join(root, "data"))
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	environmentPath := filepath.Join(root, "memory.yml")
	flagPath := filepath.Join(root, "localfs.yml")
	dataDir := filepath.Join(root, "items")
	if err := os.WriteFile(environmentPath, []byte("storage: memory\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(flagPath, []byte("storage: localfs\ndata_dir: "+dataDir+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ITEMS_EXAMPLE_CONFIG_FILE", environmentPath)

	var added bytes.Buffer
	args := []string{"items", "item", "add", "--name", "invoice", "--config", flagPath}
	if err := run(context.Background(), args, &added); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(added.String(), " invoice\n") {
		t.Fatalf("add output = %q", added.String())
	}
	var listed bytes.Buffer
	if err := run(context.Background(), []string{"items", "--config", flagPath, "item", "list"}, &listed); err != nil {
		t.Fatal(err)
	}
	if listed.String() != added.String() {
		t.Fatalf("list output = %q, want %q", listed.String(), added.String())
	}
	var memoryList bytes.Buffer
	if err := run(context.Background(), []string{"items", "item", "list"}, &memoryList); err != nil {
		t.Fatal(err)
	}
	if memoryList.Len() != 0 {
		t.Fatalf("memory config listed localfs items: %q", memoryList.String())
	}
}
