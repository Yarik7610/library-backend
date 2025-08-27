package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mail         string `mapstructure:"MAIL"`
	MailPassword string `mapstructure:"MAIL_PASSWORD"`
}

var Data Config

func Init() error {
	viper.AutomaticEnv()

	viper.BindEnv("MAIL")
	viper.BindEnv("MAIL_PASSWORD")

	err := viper.Unmarshal(&Data)
	return err
}
