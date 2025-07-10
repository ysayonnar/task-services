package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH = "config.yaml"

type RabbitMQ struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type Config struct {
	Env      string   `yaml:"env"`
	Secret   string   `yaml:"secret"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
}

func MustParse() Config {
	return MustParseByPath(CONFIG_PATH)
}

func MustParseByPath(path string) Config {
	var config Config

	configFile, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("error while reading config file: %w", err))
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(fmt.Errorf("error while parsing config: %w", err))
	}

	return config
}
