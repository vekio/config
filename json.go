package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type jsonCodec[T any] struct{}

func (jsonCodec[T]) unmarshal(content []byte, data *T) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode JSON: multiple values are not allowed")
		}
		return fmt.Errorf("decode JSON: %w", err)
	}
	return nil
}

func (jsonCodec[T]) marshal(data T) ([]byte, error) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode JSON: %w", err)
	}
	content = append(content, '\n')
	return content, nil
}
