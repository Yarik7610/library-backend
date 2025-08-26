package repository

import (
	"errors"

	"github.com/Yarik7610/library-backend/user-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(ID uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	GetEmailsByUserIDs(userIDs []uint) ([]string, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(ID uint) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetEmailsByUserIDs(userIDs []uint) ([]string, error) {
	var emails []string
	if err := r.db.Model(&model.User{}).Order("email ASC").Where("id IN ?", userIDs).Pluck("email", &emails).Error; err != nil {
		return nil, err
	}
	return emails, nil
}
