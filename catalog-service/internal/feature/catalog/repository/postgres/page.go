package postgres

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"

	postgresInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type PageRepository interface {
	WithinTX(tx *gorm.DB) PageRepository
	CreatePage(page *model.Page) error
	GetPage(bookID uint, pageNumber uint) (*model.Page, error)
}

type pageRepository struct {
	name string
	db   *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{name: "Page", db: db}
}

func (r *pageRepository) WithinTX(tx *gorm.DB) PageRepository {
	return &pageRepository{db: tx}
}

func (r *pageRepository) CreatePage(page *model.Page) error {
	if err := r.db.Create(page).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *pageRepository) GetPage(bookID, pageNumber uint) (*model.Page, error) {
	var page model.Page
	if err := r.db.Where("book_id = ?", bookID).Where("number = ?", pageNumber).First(&page).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &page, nil
}
