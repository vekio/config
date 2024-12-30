package config

import "sync"

var (
	instance *Config
	once     sync.Once
)

func Instance() *Config {
	once.Do(func() {
		instance = newConfig()
	})
	return instance
}
