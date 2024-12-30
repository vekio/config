package config

import (
	"sync"
	"testing"
)

func TestSingletonConfigInstance(t *testing.T) {
	var wg sync.WaitGroup
	instanceCount := 30
	instances := make([]*Config, instanceCount)

	for i := 0; i < instanceCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			instances[index] = Instance()
		}(i)
	}

	wg.Wait()

	firstInstance := instances[0]
	for _, instance := range instances {
		if instance != firstInstance {
			t.Errorf("Instance() did not return the same instance, got different instances")
		}
	}
}

func TestConfigFields(t *testing.T) {
	config := Instance()

	if config.appName == "" {
		t.Errorf("appName should not be empty")
	}
	if config.dir == "" {
		t.Errorf("dir should not be empty")
	}
	if config.file == "" {
		t.Errorf("file should not be empty")
	}
}
