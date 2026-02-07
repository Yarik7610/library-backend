package postgres

import (
	"context"
	"time"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"

	postgresInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type PageRepository interface {
	WithinTX(tx *gorm.DB) PageRepository
	Create(ctx context.Context, page *model.Page) error
	FindByBookIDAndPageNumber(ctx context.Context, bookID uint, pageNumber uint) (*model.Page, error)
}

type pageRepository struct {
	name    string
	timeout time.Duration
	db      *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{name: "Page(s)", timeout: 500 * time.Millisecond, db: db}
}

func (r *pageRepository) WithinTX(tx *gorm.DB) PageRepository {
	return &pageRepository{name: "Page(s)", timeout: 500 * time.Millisecond, db: tx}
}

func (r *pageRepository) Create(ctx context.Context, page *model.Page) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(page).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *pageRepository) FindByBookIDAndPageNumber(ctx context.Context, bookID, pageNumber uint) (*model.Page, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var page model.Page
	if err := r.db.WithContext(ctx).Where("book_id = ?", bookID).Where("number = ?", pageNumber).First(&page).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &page, nil
}
