package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	PostgresURL string `mapstructure:"POSTGRES_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("POSTGRES_URL")
	viper.BindEnv("JWT_SECRET")

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
