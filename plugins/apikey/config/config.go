package config

import (
	"errors"
)

type config struct {
	ApiKeyBindingName string
}

func New(cfg map[string]string) (*config, error) {
	if cfg == nil {
		return nil, errors.New("found nil config")
	}
	return buildConfig(cfg)
}

func buildConfig(cfg map[string]string) (*config, error) {
	val, ok := cfg["apiKeyBindingName"]
	if !ok {
		return nil, errors.New("required field apiKeyBindingName not found in config")
	}
	return &config{
		ApiKeyBindingName: val,
	}, nil
}
