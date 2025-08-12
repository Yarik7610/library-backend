package repository

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	CreateAuthor(author *model.Author) error
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) CreateAuthor(author *model.Author) error {
	return r.db.Create(author).Error
}
