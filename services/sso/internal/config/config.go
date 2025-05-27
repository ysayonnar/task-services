package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH = "config/config.yaml"

type GRPC struct {
	Port int `yaml:"port"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
	SslMode  string `yaml:"ssl_mode"`
}

type Config struct {
	Env      string `yaml:"env"`
	TokenTTL string `yaml:"token_ttl"`
	Secret   string `yaml:"secret"`
	GRPC     GRPC   `yaml:"grpc"`
	DB       DB     `yaml:"db"`
}

func Parse() Config {
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
