package repository

import (
	"errors"

	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"gorm.io/gorm"
)

type PageRepository interface {
	CreatePage(page *model.Page) error
	GetPage(bookID uint, pageNumber uint) (*model.Page, error)
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

func (r *pageRepository) GetPage(bookID, pageNumber uint) (*model.Page, error) {
	var page model.Page
	if err := r.db.Where("book_id = ?", bookID).Where("number = ?", pageNumber).First(&page).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &page, nil
}
