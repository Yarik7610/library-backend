package connect

import (
	"github.com/Yarik7610/library-backend/catalog-service/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.Data.PostgresURL), &gorm.Config{})
	if err != nil {
		zap.S().Fatalf("GORM open error: %v\n", err)
	}
	zap.S().Info("Successfully connected to Postgres")

	err = db.AutoMigrate(&model.Author{}, &model.Book{}, &model.Page{})
	if err != nil {
		zap.S().Fatalf("GORM auto migrate error: %v", err)
	}
	zap.S().Info("Successfully made auto migrate")

	return db
}
