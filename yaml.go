package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"go.yaml.in/yaml/v4"
)

type yamlCodec[T any] struct{}

func (yamlCodec[T]) unmarshal(content []byte, data *T) error {
	loader, err := yaml.NewLoader(bytes.NewReader(content), yaml.WithKnownFields(), yaml.WithUniqueKeys())
	if err != nil {
		return fmt.Errorf("create YAML loader: %w", err)
	}
	if err := loader.Load(data); err != nil {
		return fmt.Errorf("decode YAML: %w", err)
	}
	var extra any
	if err := loader.Load(&extra); err == nil {
		return fmt.Errorf("decode YAML: multiple documents are not allowed")
	} else if !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode YAML: %w", err)
	}
	return nil
}

func (yamlCodec[T]) marshal(data T) ([]byte, error) {
	content, err := yaml.Dump(data, yaml.WithIndent(2))
	if err != nil {
		return nil, fmt.Errorf("encode YAML: %w", err)
	}
	return content, nil
}
