package repository

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"gorm.io/gorm"
)

type BookRepository interface {
	ListCategories() ([]string, error)
	CreateBook(book *model.Book) error
	CountBooks() (int64, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) ListCategories() ([]string, error) {
	categories := make([]string, 0)
	if err := r.db.Model(&model.Book{}).Distinct().Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *bookRepository) CreateBook(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) CountBooks() (int64, error) {
	var bookCount int64
	err := r.db.Model(&model.Book{}).Count(&bookCount).Error
	return bookCount, err
}
