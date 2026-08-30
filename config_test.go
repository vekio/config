package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	config "github.com/vekio/config"
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

func TestYAMLLoadOrCreateAndLoad(t *testing.T) {
	base := t.TempDir()
	want := testConfig{Name: "api", Port: 8080}
	file, err := config.NewYAMLConfigFile[testConfig](base, "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}

	got, err := file.LoadOrCreate(want)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("LoadOrCreate() = %#v, want %#v", got, want)
	}
	if file.Path() != filepath.Join(base, "example", "config.yml") {
		t.Fatalf("Path() = %q", file.Path())
	}
	notUsed := testConfig{Name: "other", Port: 1234}
	got, err = file.LoadOrCreate(notUsed)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("LoadOrCreate() replaced existing data: got %#v, want %#v", got, want)
	}

	if err := os.WriteFile(file.Path(), []byte("name: worker\nport: 9090\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = file.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != (testConfig{Name: "worker", Port: 9090}) {
		t.Fatalf("Load() = %#v", got)
	}
}

func TestSaveUsesPrivatePermissions(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Save(testConfig{Name: "api", Port: 8080}); err != nil {
		t.Fatal(err)
	}

	dirInfo, err := os.Stat(filepath.Dir(file.Path()))
	if err != nil {
		t.Fatal(err)
	}
	if got := dirInfo.Mode().Perm(); got != 0o700 {
		t.Fatalf("directory permissions = %o, want 700", got)
	}
	fileInfo, err := os.Stat(file.Path())
	if err != nil {
		t.Fatal(err)
	}
	if got := fileInfo.Mode().Perm(); got != 0o600 {
		t.Fatalf("file permissions = %o, want 600", got)
	}
}

func TestCreateDoesNotReplaceExistingFile(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	original := testConfig{Name: "original", Port: 8080}
	if err := file.Create(original); err != nil {
		t.Fatal(err)
	}
	if err := file.Create(testConfig{Name: "replacement", Port: 9090}); !errors.Is(err, os.ErrExist) {
		t.Fatalf("Create() error = %v, want os.ErrExist", err)
	}
	loaded, err := file.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded != original {
		t.Fatalf("Load() = %#v, want %#v", loaded, original)
	}
}

func TestJSONRoundTripAndExactFilename(t *testing.T) {
	base := t.TempDir()
	file, err := config.NewJSONConfigFile[testConfig](base, "example", "settings.conf")
	if err != nil {
		t.Fatal(err)
	}

	want := testConfig{Name: "api", Port: 443}
	if err := file.Save(want); err != nil {
		t.Fatal(err)
	}
	if filepath.Base(file.Path()) != "settings.conf" {
		t.Fatalf("filename = %q", filepath.Base(file.Path()))
	}

	second, err := config.NewJSONConfigFile[testConfig](base, "example", "settings.conf")
	if err != nil {
		t.Fatal(err)
	}
	got, err := second.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
}

func TestYAMLAcceptsYMLFilename(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "settings.yml")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(file.Path()) != "settings.yml" {
		t.Fatalf("filename = %q, want settings.yml", filepath.Base(file.Path()))
	}
}

func TestLoadRejectsInvalidData(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](
		t.TempDir(),
		"example",
		"config.yml",
	)
	if err != nil {
		t.Fatal(err)
	}
	original := testConfig{Name: "valid", Port: 80}
	if err := file.Save(original); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file.Path(), []byte("name: ''\nport: 80\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err = file.Load()
	if err == nil || !strings.Contains(err.Error(), "validate configuration") {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestStrictDecoding(t *testing.T) {
	tests := []struct {
		name    string
		newFile func(string) (*config.ConfigFile[testConfig], error)
		content string
	}{
		{
			name: "yaml unknown field",
			newFile: func(path string) (*config.ConfigFile[testConfig], error) {
				return config.NewYAMLConfigFile[testConfig](path, "app", "config.yml")
			},
			content: "name: api\nport: 80\nunknown: true\n",
		},
		{
			name: "json unknown field",
			newFile: func(path string) (*config.ConfigFile[testConfig], error) {
				return config.NewJSONConfigFile[testConfig](path, "app", "config.json")
			},
			content: `{"name":"api","port":80,"unknown":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := tt.newFile(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Dir(file.Path()), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file.Path(), []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := file.Load(); err == nil {
				t.Fatal("Load() succeeded for an unknown field")
			}
		})
	}
}

func TestConstructorRejectsUnsafeNames(t *testing.T) {
	_, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "../escape", "config.yml")
	if err == nil {
		t.Fatal("constructor accepted an unsafe application name")
	}
}

func TestConstructorRejectsEmptyBaseDirectory(t *testing.T) {
	_, err := config.NewYAMLConfigFile[testConfig]("", "example", "config.yml")
	if err == nil {
		t.Fatal("constructor accepted an empty base directory")
	}
}

func TestSetPathOverridesConventionalPath(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "my-app", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(t.TempDir(), "development.yml")
	if err := file.SetPath(want); err != nil {
		t.Fatal(err)
	}
	if file.AppName() != "my-app" {
		t.Fatalf("AppName() = %q, want my-app", file.AppName())
	}
	if file.Path() != want {
		t.Fatalf("Path() = %q, want %q", file.Path(), want)
	}
}

func TestSetPathRejectsEmptyPath(t *testing.T) {
	file, err := config.NewYAMLConfigFile[testConfig](t.TempDir(), "example", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.SetPath(" "); err == nil {
		t.Fatal("SetPath() accepted an empty path")
	}
}

func TestSaveRejectsFileAsBaseDirectory(t *testing.T) {
	baseFile := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(baseFile, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := config.NewYAMLConfigFile[testConfig](baseFile, "example", "config.yml")
	if err != nil {
		t.Fatalf("constructor performed filesystem validation: %v", err)
	}
	if err := file.Save(testConfig{Name: "example", Port: 80}); err == nil {
		t.Fatal("Save() accepted a file as its base directory")
	}
}

func TestDefaultConfigFilePath(t *testing.T) {
	configHome := filepath.Join(t.TempDir(), "config-home")
	t.Setenv("XDG_CONFIG_HOME", configHome)
	file, err := config.NewDefaultConfigFile[testConfig]("example")
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(configHome, "example", "config.yml")
	if file.Path() != want {
		t.Fatalf("Path() = %q, want %q", file.Path(), want)
	}
}

func TestDefaultDataDir(t *testing.T) {
	t.Run("XDG_DATA_HOME", func(t *testing.T) {
		base := filepath.Join(t.TempDir(), "data")
		t.Setenv("XDG_DATA_HOME", base)

		got, err := config.DefaultDataDir("example")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(base, "example")
		if got != want {
			t.Fatalf("DefaultDataDir() = %q, want %q", got, want)
		}
	})

	t.Run("relative XDG_DATA_HOME", func(t *testing.T) {
		t.Setenv("XDG_DATA_HOME", "relative/data")
		if _, err := config.DefaultDataDir("example"); err == nil {
			t.Fatal("DefaultDataDir() accepted a relative XDG_DATA_HOME")
		}
	})

	t.Run("home fallback", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("XDG_DATA_HOME", "")
		t.Setenv("HOME", home)

		got, err := config.DefaultDataDir("example")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, ".local", "share", "example")
		if got != want {
			t.Fatalf("DefaultDataDir() = %q, want %q", got, want)
		}
	})

	t.Run("unsafe application name", func(t *testing.T) {
		if _, err := config.DefaultDataDir("../escape"); err == nil {
			t.Fatal("DefaultDataDir() accepted an unsafe application name")
		}
	})
}
