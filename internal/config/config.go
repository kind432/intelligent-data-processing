package config

import (
	"github.com/spf13/viper"
)

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AddConfigPath(viper.GetString("env.path"))

	err = viper.MergeInConfig()
	if err != nil {
		return err
	}
	return nil
}
