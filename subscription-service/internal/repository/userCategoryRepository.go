package repository

import (
	"gorm.io/gorm"
)

type UserCategoryRepository interface {
	GetUserCategoriesSubscriptions(userID uint) ([]string, error)
}

type userCategoryRepository struct {
	db *gorm.DB
}

func NewUserCategoryRepository(db *gorm.DB) UserCategoryRepository {
	return &userCategoryRepository{db: db}
}

func (r *userCategoryRepository) GetUserCategoriesSubscriptions(userID uint) ([]string, error) {
	return nil, nil
}
