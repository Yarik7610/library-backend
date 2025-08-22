package repository

import (
	"errors"

	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	CreateAuthor(author *model.Author) error
	FindByID(ID uint) (*model.Author, error)
	DeleteAuthor(ID uint) error
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

func (r *authorRepository) FindByID(ID uint) (*model.Author, error) {
	var author model.Author
	if err := r.db.Where("id = ?", ID).First(&author).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &author, nil
}

func (r *authorRepository) DeleteAuthor(ID uint) error {
	return r.db.Delete(&model.Author{}, ID).Error
}
