package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MemoryRepository struct{ items []Item }

func (r *MemoryRepository) Save(item Item) error {
	r.items = append(r.items, item)
	return nil
}

func (r *MemoryRepository) List() ([]Item, error) {
	return append([]Item(nil), r.items...), nil
}

type LocalFSRepository struct{ Dir string }

func (r *LocalFSRepository) Save(item Item) error {
	if err := os.MkdirAll(r.Dir, 0o700); err != nil {
		return fmt.Errorf("create item directory: %w", err)
	}
	content, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("encode item: %w", err)
	}
	path := filepath.Join(r.Dir, item.ID+".json")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create item: %w", err)
	}
	if _, err := file.Write(append(content, '\n')); err != nil {
		file.Close()
		os.Remove(path)
		return fmt.Errorf("write item: %w", err)
	}
	if err := file.Close(); err != nil {
		os.Remove(path)
		return fmt.Errorf("close item: %w", err)
	}
	return nil
}

func (r *LocalFSRepository) List() ([]Item, error) {
	entries, err := os.ReadDir(r.Dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read item directory: %w", err)
	}
	var items []Item
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(r.Dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read item %s: %w", entry.Name(), err)
		}
		var item Item
		if err := json.Unmarshal(content, &item); err != nil {
			return nil, fmt.Errorf("decode item %s: %w", entry.Name(), err)
		}
		items = append(items, item)
	}
	return items, nil
}
