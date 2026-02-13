package postgres

import (
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres/model"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.PostgresURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&model.UserBookCategory{}); err != nil {
		return nil, err
	}

	return db, nil
}
