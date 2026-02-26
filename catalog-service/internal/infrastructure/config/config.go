package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                      string `env:"ENV"`
	ServiceName              string `env:"SERVICE_NAME"`
	HTTPServerPort           string `env:"HTTP_SERVER_PORT"`
	GRPCServerPort           string `env:"GRPC_SERVER_PORT"`
	PostgresURL              string `env:"POSTGRES_URL"`
	RedisHost                string `env:"REDIS_HOST"`
	RedisPort                string `env:"REDIS_PORT"`
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
