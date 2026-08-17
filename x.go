package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultDirMode  = 0o700
	defaultFileMode = 0o600
)

func ensureDir(path string) error {
	return os.MkdirAll(path, defaultDirMode)
}

func normalizeDirPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("directory path cannot be empty")
	}
	path = filepath.Clean(path)
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return "", fmt.Errorf("%q is not a directory", path)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	return path, nil
}

func normalizeFilePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("file path cannot be empty")
	}
	path = filepath.Clean(path)

	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return "", fmt.Errorf("%q is a directory", path)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	parentInfo, parentErr := os.Stat(filepath.Dir(path))
	if parentErr == nil && !parentInfo.IsDir() {
		return "", fmt.Errorf("parent of %q is not a directory", path)
	}
	if parentErr != nil && !errors.Is(parentErr, os.ErrNotExist) {
		return "", parentErr
	}
	return path, nil
}

func writeFileAtomic(path string, content []byte) (resultErr error) {
	dir := filepath.Dir(path)
	if err := ensureDir(dir); err != nil {
		return err
	}

	temporary, err := os.CreateTemp(dir, ".config-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		if resultErr != nil {
			_ = os.Remove(temporaryPath)
		}
	}()

	if err := temporary.Chmod(defaultFileMode); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(content); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

func writeFileExclusive(path string, content []byte) (resultErr error) {
	dir := filepath.Dir(path)
	if err := ensureDir(dir); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, defaultFileMode)
	if err != nil {
		return err
	}
	defer func() {
		if resultErr != nil {
			_ = file.Close()
			_ = os.Remove(path)
		}
	}()

	written, err := file.Write(content)
	if err != nil {
		return err
	}
	if written != len(content) {
		return io.ErrShortWrite
	}
	if err := file.Sync(); err != nil {
		return err
	}
	return file.Close()
}
