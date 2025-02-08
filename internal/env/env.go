package env

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      App      `yaml:"app"`
	Telegram Telegram `yaml:"telegram"`
	Storage  Storage  `yaml:"storage"`
}

type App struct {
	Name  string `yaml:"name" env:"APP_NAME" env-default:"tBot"`
	Debug bool   `yaml:"debug" env:"APP_DEBUG" env-default:"false"`
}
type Telegram struct {
	Host  string `yaml:"host" env:"TELEGRAM_HOST" env-default:"api.telegram.org"`
	Token string `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
}

type Storage struct {
	DB string `yaml:"db" env:"STORAGE_DB" env-default:"./data/bot.db"`
}

func MustLoadConfig() *Config {
	var (
		cfg Config
		err error
	)
	_ = cleanenv.ReadEnv(&cfg)
	files := []string{
		"./config.yaml",
		"./.env",
	}
	for _, file := range files {
		if err = cleanenv.ReadConfig(file, &cfg); err != nil {
			continue
		}
	}
	return &cfg
}
