package postgres

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.PostgresURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Author{}, &model.Book{}, &model.Page{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
