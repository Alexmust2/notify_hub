package config

import "os"

type Config struct {
    TelegramBotToken string
}

func Load() *Config {
    return &Config{
        TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
    }
}
