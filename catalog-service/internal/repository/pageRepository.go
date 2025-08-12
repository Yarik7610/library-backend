package repository

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"gorm.io/gorm"
)

type PageRepository interface {
	CreatePage(page *model.Page) error
}

type pageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}

func (r *pageRepository) CreatePage(page *model.Page) error {
	return r.db.Create(page).Error
}
