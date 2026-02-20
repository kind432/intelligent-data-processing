package config

import (
	"github.com/spf13/viper"
)

func New() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetConfigName("development")
	viper.SetConfigType("env")
	viper.AddConfigPath(viper.GetString("env.path"))
	_ = viper.MergeInConfig() // Игнорируем ошибку, если dev.env нет
	return nil
}
