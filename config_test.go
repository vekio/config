package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testSettings struct {
	Name string `json:"name" yaml:"name"`
	Port int    `json:"port" yaml:"port"`
}

func (t testSettings) Validate() error { return nil }

func TestNewYAMLConfigFile(t *testing.T) {
	cfg, err := NewYAMLConfigFile[testSettings]()
	if err != nil {
		t.Fatalf("NewYAMLConfigFile returned error: %v", err)
	}

	if cfg == nil {
		t.Fatalf("expected config file, got nil")
	}

	if cfg.fileName != "config.yml" {
		t.Fatalf("expected file name config.yml, got %q", cfg.fileName)
	}

	if ext := cfg.fileManager.Extension(); ext != ".yml" {
		t.Fatalf("expected extension .yml, got %q", ext)
	}
}

func TestNewJSONConfigFile(t *testing.T) {
	cfg, err := NewJSONConfigFile[testSettings]()
	if err != nil {
		t.Fatalf("NewJSONConfigFile returned error: %v", err)
	}

	if cfg.fileName != "config.json" {
		t.Fatalf("expected file name config.json, got %q", cfg.fileName)
	}

	if ext := cfg.fileManager.Extension(); ext != ".json" {
		t.Fatalf("expected extension .json, got %q", ext)
	}
}

func TestNewConfigFileNilManager(t *testing.T) {
	if _, err := newConfigFile[testSettings](nil); err == nil {
		t.Fatalf("expected error when manager is nil")
	}
}

func TestWithDefault(t *testing.T) {
	cfg := newTestConfigFile(t)
	defaults := testSettings{Name: "default", Port: 1}

	WithDefault(defaults)(cfg)

	if cfg.defaultData != defaults {
		t.Fatalf("expected default data %+v, got %+v", defaults, cfg.defaultData)
	}
}

func TestWithDefaultIgnoresNilConfig(t *testing.T) {
	defaults := testSettings{Name: "default", Port: 1}
	WithDefault(defaults)(nil)
}

func TestWithAppName(t *testing.T) {
	cfg := newTestConfigFile(t)
	cfg.appName = "original"

	WithAppName[testSettings]("  demo  ")(cfg)

	if cfg.appName != "demo" {
		t.Fatalf("expected app name demo, got %q", cfg.appName)
	}

	WithAppName[testSettings]("   ")(cfg)
	if cfg.appName != "demo" {
		t.Fatalf("expected app name to remain demo, got %q", cfg.appName)
	}
}

func TestWithPath(t *testing.T) {
	cfg := newTestConfigFile(t)
	original := cfg.path
	WithPath[testSettings]("  /tmp/path  ")(cfg)
	if !strings.HasSuffix(cfg.path, "/tmp/path") {
		t.Fatalf("expected path to end with /tmp/path, got %q", cfg.path)
	}

	WithPath[testSettings]("   ")(cfg)
	if cfg.path != strings.TrimSpace("  /tmp/path  ") && cfg.path != original {
		t.Fatalf("expected path unchanged on empty input, got %q", cfg.path)
	}
}

func TestWithFilename(t *testing.T) {
	cfg := newTestConfigFile(t)

	WithFilename[testSettings]("settings")(cfg)
	if cfg.fileName != "settings.json" {
		t.Fatalf("expected file name settings.json, got %q", cfg.fileName)
	}

	WithFilename[testSettings]("custom.yaml")(cfg)
	if cfg.fileName != "custom.json" {
		t.Fatalf("expected file name custom.json, got %q", cfg.fileName)
	}
}

func TestDirPath(t *testing.T) {
	cfg := mustNewTestConfigFile(t, WithAppName[testSettings]("myapp"))
	expected := filepath.Join(cfg.path, "myapp")
	if got := cfg.DirPath(); got != expected {
		t.Fatalf("expected dir path %q, got %q", expected, got)
	}
}

func TestPathNoEnv(t *testing.T) {
	cfg := mustNewTestConfigFile(t, WithAppName[testSettings]("myapp"))
	expected := filepath.Join(cfg.DirPath(), cfg.fileName)
	if got := cfg.Path(); got != expected {
		t.Fatalf("expected path %q, got %q", expected, got)
	}
}

func TestPathWithEnvFilePresent(t *testing.T) {
	cfg := mustNewTestConfigFile(t, WithAppName[testSettings]("myapp"))
	envVar := strings.ToUpper(cfg.appName) + "_ENV"
	t.Setenv(envVar, "DEV")

	dir := cfg.DirPath()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create dir: %v", err)
	}

	base := strings.TrimSuffix(cfg.fileName, filepath.Ext(cfg.fileName))
	rawExt := filepath.Ext(cfg.fileName)
	extension := strings.TrimPrefix(rawExt, ".")
	lowerEnv := strings.ToLower("DEV")
	var candidate string
	if extension == "" {
		candidate = fmt.Sprintf("%s.%s", base, lowerEnv)
	} else {
		candidate = fmt.Sprintf("%s.%s.%s", base, lowerEnv, extension)
	}
	candidatePath := filepath.Join(dir, candidate)
	if err := os.WriteFile(candidatePath, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write candidate file: %v", err)
	}

	if got := cfg.Path(); got != candidatePath {
		t.Fatalf("expected env-specific path %q, got %q", candidatePath, got)
	}
}

func TestContent(t *testing.T) {
	cfg := mustNewTestConfigFile(t)
	data := testSettings{Name: "content", Port: 5}
	if err := cfg.Init(data); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	content, err := cfg.Content()
	if err != nil {
		t.Fatalf("Content failed: %v", err)
	}

	want, err := os.ReadFile(cfg.Path())
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	if string(content) != string(want) {
		t.Fatalf("content mismatch: got %q want %q", string(content), string(want))
	}
}

func TestData(t *testing.T) {
	cfg := mustNewTestConfigFile(t)
	data := testSettings{Name: "data", Port: 2}
	if err := cfg.Init(data); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if got := cfg.Data(); got != data {
		t.Fatalf("expected data %+v, got %+v", data, got)
	}
}

func TestReload(t *testing.T) {
	cfg := mustNewTestConfigFile(t)
	initial := testSettings{Name: "before", Port: 10}
	if err := cfg.Init(initial); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	updated := testSettings{Name: "after", Port: 20}
	buf, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(cfg.Path(), buf, 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := cfg.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	if got := cfg.Data(); got != updated {
		t.Fatalf("expected data %+v, got %+v", updated, got)
	}
}

func TestInitCreatesFile(t *testing.T) {
	cfg := mustNewTestConfigFile(t)
	data := testSettings{Name: "init", Port: 99}
	if err := cfg.Init(data); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if _, err := os.Stat(cfg.Path()); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestSoftInitCreatesDefaults(t *testing.T) {
	defaults := testSettings{Name: "default", Port: 33}
	cfg := mustNewTestConfigFile(t, WithDefault(defaults))

	if err := cfg.SoftInit(); err != nil {
		t.Fatalf("SoftInit failed: %v", err)
	}

	if got := cfg.Data(); got != defaults {
		t.Fatalf("expected data %+v, got %+v", defaults, got)
	}

	if _, err := os.Stat(cfg.Path()); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestSoftInitLoadsExisting(t *testing.T) {
	defaults := testSettings{Name: "default", Port: 33}
	cfg := mustNewTestConfigFile(t, WithDefault(defaults))

	existing := testSettings{Name: "existing", Port: 44}
	if err := cfg.Init(existing); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	cfg.defaultData = testSettings{Name: "changed", Port: 55}

	if err := cfg.SoftInit(); err != nil {
		t.Fatalf("SoftInit failed: %v", err)
	}

	if got := cfg.Data(); got != existing {
		t.Fatalf("expected data %+v, got %+v", existing, got)
	}
}

func TestGetFileNameForEnvironment(t *testing.T) {
	dir := t.TempDir()
	appName := "demo"
	fileName := "config.json"

	t.Setenv(strings.ToUpper(appName)+"_ENV", "DEV")

	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	rawExt := filepath.Ext(fileName)
	extension := strings.TrimPrefix(rawExt, ".")
	lowerEnv := strings.ToLower("DEV")
	var candidate string
	if extension == "" {
		candidate = fmt.Sprintf("%s.%s", base, lowerEnv)
	} else {
		candidate = fmt.Sprintf("%s.%s.%s", base, lowerEnv, extension)
	}
	candidatePath := filepath.Join(dir, candidate)
	if err := os.WriteFile(candidatePath, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if got := getFileNameForEnvironment(dir, appName, fileName); got != candidate {
		t.Fatalf("expected %q, got %q", candidate, got)
	}
}

func TestJSONFileManagerRoundTrip(t *testing.T) {
	mgr := NewJSONFileManager[testSettings]()
	path := filepath.Join(t.TempDir(), "config.json")
	input := testSettings{Name: "round", Port: 7}

	if err := mgr.WriteDataToFile(path, input); err != nil {
		t.Fatalf("WriteDataToFile failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	var output testSettings
	if err := mgr.LoadDataFromFile(path, &output); err != nil {
		t.Fatalf("LoadDataFromFile failed: %v", err)
	}

	if output != input {
		t.Fatalf("expected %+v, got %+v", input, output)
	}

	if ext := mgr.Extension(); ext != ".json" {
		t.Fatalf("expected extension .json, got %q", ext)
	}
}

func TestYAMLFileManagerRoundTrip(t *testing.T) {
	mgr := NewYAMLFileManager[testSettings]()
	path := filepath.Join(t.TempDir(), "config.yml")
	input := testSettings{Name: "round", Port: 8}

	if err := mgr.WriteDataToFile(path, input); err != nil {
		t.Fatalf("WriteDataToFile failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	var output testSettings
	if err := mgr.LoadDataFromFile(path, &output); err != nil {
		t.Fatalf("LoadDataFromFile failed: %v", err)
	}

	if output != input {
		t.Fatalf("expected %+v, got %+v", input, output)
	}

	if ext := mgr.Extension(); ext != ".yml" {
		t.Fatalf("expected extension .yml, got %q", ext)
	}
}

func newTestConfigFile(t *testing.T) *ConfigFile[testSettings] {
	t.Helper()
	return &ConfigFile[testSettings]{
		fileManager: NewJSONFileManager[testSettings](),
		fileName:    "config.json",
		appName:     "app",
		path:        t.TempDir(),
	}
}

func mustNewTestConfigFile(t *testing.T, opts ...ConfigFileOption[testSettings]) *ConfigFile[testSettings] {
	t.Helper()
	tempPath := t.TempDir()
	options := []ConfigFileOption[testSettings]{
		WithPath[testSettings](tempPath),
		WithAppName[testSettings]("testapp"),
	}
	options = append(options, opts...)

	cfg, err := NewJSONConfigFile(options...)
	if err != nil {
		t.Fatalf("NewJSONConfigFile failed: %v", err)
	}
	return cfg
}
