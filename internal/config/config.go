package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram map[string]TelegramConfig `yaml:"telegram"`
	Email    map[string]EmailConfig    `yaml:"email"`
}

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}