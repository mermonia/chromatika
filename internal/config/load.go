package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

const ConfigFileName string = "chromatika.toml"

func LoadConfig(path string) (*Config, error) {
	c := &Config{}
	if _, err := toml.DecodeFile(path, c); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("config has an invalid configuration: %w", err)
	}

	return c, nil
}

func (c *Config) validate() error {
	return nil
}
