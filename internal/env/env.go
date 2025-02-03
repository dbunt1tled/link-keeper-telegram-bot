package env

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Telegram Telegram `yaml:"telegram"`
	Storage  Storage  `yaml:"storage"`
}
type Telegram struct {
	Host  string `yaml:"host" env:"TELEGRAM_HOST" env-default:"api.telegram.org"`
	Token string `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
}

type Storage struct {
	BasePath string `yaml:"base_path" env:"STORAGE_BASE_PATH" env-default:"data"`
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
