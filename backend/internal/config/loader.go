package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func LoadConfig() (*ServiceConfig, error) {
	k := koanf.New(".")

	// Load configuration from file (config.yaml)
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		return nil, err
	}

	// Load environment variables and override the values from the file
	if err := k.Load(env.Provider("", "_", func(s string) string {
		return strings.ToLower(s)
	}), nil); err != nil {
		return nil, err
	}

	var cfg ServiceConfig
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
