package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mail         string `mapstructure:"MAIL"`
	MailPassword string `mapstructure:"MAIL_PASSWORD"`
	Env          string `mapstructure:"ENV"`
}

func Init() (*Config, error) {
	viper.AutomaticEnv()

	viper.BindEnv("MAIL")
	viper.BindEnv("MAIL_PASSWORD")
	viper.BindEnv("ENV")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
