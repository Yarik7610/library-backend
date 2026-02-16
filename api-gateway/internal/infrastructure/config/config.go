package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	Env        string `mapstructure:"ENV"`
}

func Init() (*Config, error) {
	viper.AutomaticEnv()

	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("ENV")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
