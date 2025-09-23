package config

import (
	"fmt"
	"sync"
)

var (
	configSingletonMu sync.RWMutex
	configSingleton   any
)

// SetConfigFile registers the provided ConfigFile as the global instance used by helper packages like cli.
func SetConfigFile[T Validatable](cfg *ConfigFile[T]) {
	configSingletonMu.Lock()
	defer configSingletonMu.Unlock()

	configSingleton = cfg
}

// ConfigFileSingleton returns the registered ConfigFile, or an error when it has not yet been configured.
func ConfigFileSingleton[T Validatable]() (*ConfigFile[T], error) {
	configSingletonMu.RLock()
	defer configSingletonMu.RUnlock()

	if configSingleton == nil {
		return nil, fmt.Errorf("config: configuration file not registered; call SoftInit or Init first")
	}

	cfg, ok := configSingleton.(*ConfigFile[T])
	if !ok {
		return nil, fmt.Errorf("config: registered configuration file has a different type parameter")
	}

	return cfg, nil
}

// MustConfigFile returns the registered ConfigFile, panicking if it has not been registered yet.
func MustConfigFile[T Validatable]() *ConfigFile[T] {
	cfg, err := ConfigFileSingleton[T]()
	if err != nil {
		panic(err)
	}
	return cfg
}
