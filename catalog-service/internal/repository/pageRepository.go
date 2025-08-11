package repository

import (
	"gorm.io/gorm"
)

type PageRepository interface {
}

type pageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}
