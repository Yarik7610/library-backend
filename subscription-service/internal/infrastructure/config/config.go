package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                      string `env:"ENV"`
	ServiceName              string `env:"SERVICE_NAME"`
	HTTPServerPort           string `env:"HTTP_SERVER_PORT"`
	PostgresURL              string `env:"POSTGRES_URL"`
	JWTSecret                string `env:"JWT_SECRET"`
	OTelExporterOTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
}

func Parse() (*Config, error) {
	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
