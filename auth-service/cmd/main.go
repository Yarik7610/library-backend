package main

import (
	"auth-service/config"
	"auth-service/internal/model"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var logger *zap.SugaredLogger

func main() {
	baseLogger, _ := zap.NewDevelopment()
	defer baseLogger.Sync()
	logger = baseLogger.Sugar()

	config, err := config.Load()
	if err != nil {
		logger.Fatalf("Config load error: %v\n", err)
	}

	db, err := gorm.Open(postgres.Open(config.PostgresURL), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Gorm open error: %v\n", err)
	}
	logger.Info("Successfully connected to Postgres")

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		logger.Fatalf("Gorm auto migrate error: %v", err)
	}
	logger.Info("Successfully made auto migrate")
}
