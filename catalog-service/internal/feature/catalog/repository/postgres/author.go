package postgres

import (
	"context"
	"time"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
	postgresInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	WithinTX(tx *gorm.DB) AuthorRepository
	Create(ctx context.Context, author *model.Author) error
	FindByID(ctx context.Context, authorID uint) (*model.Author, error)
	Delete(ctx context.Context, authorID uint) error
}

type authorRepository struct {
	name    string
	timeout time.Duration
	db      *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{name: "Author(s)", timeout: 500 * time.Millisecond, db: db}
}

func (r *authorRepository) WithinTX(tx *gorm.DB) AuthorRepository {
	return &authorRepository{name: "Author(s)", timeout: 500 * time.Millisecond, db: tx}
}

func (r *authorRepository) Create(ctx context.Context, author *model.Author) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(author).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *authorRepository) FindByID(ctx context.Context, authorID uint) (*model.Author, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var author model.Author
	if err := r.db.WithContext(ctx).Where("id = ?", authorID).First(&author).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &author, nil
}

func (r *authorRepository) Delete(ctx context.Context, authorID uint) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.WithContext(ctx).Delete(&model.Author{}, authorID).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}
