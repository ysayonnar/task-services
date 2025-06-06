package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH = "config.yaml"

// TODO: дописать поля для сетапа http-сервера
type Config struct {
	Env string `yaml:"env"`
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
