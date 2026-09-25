package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

type Config struct {
	Storage string `yaml:"storage"`
	DataDir string `yaml:"data_dir"`
}

func (c Config) Validate() error {
	switch c.Storage {
	case "memory":
		return nil
	case "localfs":
		if strings.TrimSpace(c.DataDir) == "" {
			return fmt.Errorf("data_dir is required for localfs storage")
		}
		return nil
	default:
		return fmt.Errorf("storage must be memory or localfs")
	}
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ItemRepository interface {
	Save(Item) error
	List() ([]Item, error)
}

type App struct {
	AddItem   AddItem
	ListItems ListItems
}

func NewApp(cfg Config) (*App, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	var repository ItemRepository
	switch cfg.Storage {
	case "memory":
		repository = &MemoryRepository{}
	case "localfs":
		repository = &LocalFSRepository{Dir: cfg.DataDir}
	}
	return &App{
		AddItem:   AddItem{Repository: repository},
		ListItems: ListItems{Repository: repository},
	}, nil
}

type AddItem struct{ Repository ItemRepository }

func (uc AddItem) Run(name string) (Item, error) {
	if strings.TrimSpace(name) == "" {
		return Item{}, fmt.Errorf("item name is required")
	}
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return Item{}, fmt.Errorf("generate item ID: %w", err)
	}
	item := Item{ID: hex.EncodeToString(id[:]), Name: name}
	if err := uc.Repository.Save(item); err != nil {
		return Item{}, err
	}
	return item, nil
}

type ListItems struct{ Repository ItemRepository }

func (uc ListItems) Run() ([]Item, error) { return uc.Repository.List() }
