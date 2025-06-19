package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH = "config.yaml"

type Config struct {
	Env          string `yaml:"env"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	Secret       string `yaml:"secret"`
}

func MustParse() Config {
	var config Config

	configFile, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		panic(fmt.Errorf("error while reading config file: %w", err))
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(fmt.Errorf("error while parsing config: %w", err))
	}

	return config
}
