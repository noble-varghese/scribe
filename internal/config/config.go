package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Expansions map[string]string `yaml:"expansions"`
}

func Load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
