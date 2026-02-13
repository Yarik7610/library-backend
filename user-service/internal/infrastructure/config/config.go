package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort           string `mapstructure:"SERVER_PORT"`
	PostgresURL          string `mapstructure:"POSTGRES_URL"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	JWTExpirationSeconds uint   `mapstructure:"JWT_EXPIRATION_SECONDS"`
	Mail                 string `mapstructure:"MAIL"`
	Env                  string `mapstructure:"ENV"`
}

func Init() (*Config, error) {
	viper.AutomaticEnv()

	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("POSTGRES_URL")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_EXPIRATION_SECONDS")
	viper.BindEnv("MAIL")
	viper.BindEnv("ENV")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
