package config

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"
)

type jsonCodec[T any] struct{}

func (jsonCodec[T]) unmarshal(content []byte, data *T) error {
	if err := jsonv2.Unmarshal(
		content,
		data,
		jsonv2.RejectUnknownMembers(true),
		jsontext.AllowDuplicateNames(false),
	); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	return nil
}

func (jsonCodec[T]) marshal(data T) ([]byte, error) {
	content, err := jsonv2.Marshal(
		data,
		jsontext.WithIndent("  "),
	)
	if err != nil {
		return nil, fmt.Errorf("encode JSON: %w", err)
	}
	content = append(content, '\n')
	return content, nil
}
