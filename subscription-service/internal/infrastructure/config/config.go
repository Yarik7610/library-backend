package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	PostgresURL string `mapstructure:"POSTGRES_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	RedisHost   string `mapstructure:"REDIS_HOST"`
	RedisPort   string `mapstructure:"REDIS_PORT"`
}

var Data Config

func Init() error {
	viper.AutomaticEnv()

	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("POSTGRES_URL")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")

	err := viper.Unmarshal(&Data)
	return err
}
