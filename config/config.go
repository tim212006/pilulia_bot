package config

import (
	"errors"
	"github.com/spf13/viper"
	"pilulia_bot/logger/consts"
)

type Config struct {
	BotSettings BotSettings
	DBSettings  DBSettings
}

type BotSettings struct {
	TelegramToken string
}

type DBSettings struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func ConfigInit() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New(consts.ErrorEnv)
	}
	config := &Config{
		BotSettings: BotSettings{
			TelegramToken: viper.GetString("TELEGRAM_TOKEN"),
		},
		DBSettings: DBSettings{
			Username: viper.GetString("USER"),
			Password: viper.GetString("PASSWORD"),
			Host:     viper.GetString("HOST"),
			Port:     viper.GetString("PORT"),
			Database: viper.GetString("DATABASE"),
		},
	}
	return config, nil
}
