package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH = "config.yaml"

type GRPC struct {
	Port int `yaml:"port"`
}

type RabbitMQ struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type Config struct {
	Env      string   `yaml:"env"`
	TokenTTL string   `yaml:"token_ttl"`
	Secret   string   `yaml:"secret"`
	GRPC     GRPC     `yaml:"grpc"`
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
