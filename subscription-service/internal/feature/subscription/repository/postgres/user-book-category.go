package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres/model"

	postgresInfrastructure "github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type UserBookCategoryRepository interface {
	FindUserBookCategory(ctx context.Context, userID uint, category string) (*model.UserBookCategory, error)
	GetBookCategoryUserIDs(ctx context.Context, category string) ([]uint, error)
	GetUserBookCategories(ctx context.Context, userID uint) ([]string, error)
	Create(ctx context.Context, userBookCategory *model.UserBookCategory) error
	Delete(ctx context.Context, userID uint, category string) error
}

type userBookCategoryRepository struct {
	name    string
	timeout time.Duration
	db      *gorm.DB
}

func NewUserBookCategoryRepository(db *gorm.DB) UserBookCategoryRepository {
	return &userBookCategoryRepository{
		name:    "Book category subscription(s)",
		timeout: 500 * time.Millisecond,
		db:      db,
	}
}

func (r *userBookCategoryRepository) FindUserBookCategory(ctx context.Context, userID uint, category string) (*model.UserBookCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var userBookCategory model.UserBookCategory
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("book_category ILIKE ?", category).
		First(&userBookCategory).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &userBookCategory, nil
}

func (r *userBookCategoryRepository) GetBookCategoryUserIDs(ctx context.Context, category string) ([]uint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var emails []uint
	if err := r.db.WithContext(ctx).
		Model(&model.UserBookCategory{}).
		Order("user_id ASC").
		Where("book_category ILIKE ?", category).
		Pluck("user_id", &emails).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return emails, nil
}

func (r *userBookCategoryRepository) GetUserBookCategories(ctx context.Context, userID uint) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var subscribedCategories []string
	if err := r.db.WithContext(ctx).
		Model(&model.UserBookCategory{}).
		Order("created_at DESC").
		Where("user_id = ?", userID).
		Pluck("book_category", &subscribedCategories).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return subscribedCategories, nil
}

func (r *userBookCategoryRepository) Create(ctx context.Context, userBookCategory *model.UserBookCategory) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	userBookCategory.BookCategory = strings.ToLower(userBookCategory.BookCategory)

	if err := r.db.WithContext(ctx).Create(userBookCategory).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *userBookCategoryRepository) Delete(ctx context.Context, userID uint, category string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("book_category ILIKE ?", category).
		Delete(&model.UserBookCategory{}).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}
