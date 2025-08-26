package repository

import (
	"errors"
	"strings"

	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserCategoryRepository interface {
	FindSubscribedCategory(userID uint, category string) (*model.UserCategory, error)
	GetCategorySubscribersIDs(category string) ([]uint, error)
	GetSubscribedCategories(userID uint) ([]string, error)
	SubscribeCategory(userCategory *model.UserCategory) error
	UnsubscribeCategory(userID uint, category string) error
}

type userCategoryRepository struct {
	db *gorm.DB
}

func NewUserCategoryRepository(db *gorm.DB) UserCategoryRepository {
	return &userCategoryRepository{db: db}
}

func (r *userCategoryRepository) FindSubscribedCategory(userID uint, category string) (*model.UserCategory, error) {
	var subscribedCategory model.UserCategory
	if err := r.db.Where("user_id = ?", userID).Where("category ILIKE ?", category).First(&subscribedCategory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &subscribedCategory, nil
}

func (r *userCategoryRepository) GetCategorySubscribersIDs(category string) ([]uint, error) {
	var emails []uint
	if err := r.db.Model(&model.UserCategory{}).Order("user_id ASC").Where("category ILIKE ?", category).Pluck("user_id", &emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}

func (r *userCategoryRepository) GetSubscribedCategories(userID uint) ([]string, error) {
	var subscribedCategories []string
	if err := r.db.Model(&model.UserCategory{}).Order("created_at DESC").Where("user_id = ?", userID).Pluck("category", &subscribedCategories).Error; err != nil {
		return nil, err
	}
	return subscribedCategories, nil
}

func (r *userCategoryRepository) SubscribeCategory(userCategory *model.UserCategory) error {
	userCategory.Category = strings.ToLower(userCategory.Category)

	err := r.db.Create(userCategory).Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return errors.New("can't subscribe on subscribed category twice")
	}
	return err
}

func (r *userCategoryRepository) UnsubscribeCategory(userID uint, category string) error {
	return r.db.Where("user_id = ?", userID).Where("category ILIKE ?", category).Delete(&model.UserCategory{}).Error
}
