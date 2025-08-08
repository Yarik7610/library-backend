package service

import (
	"errors"
	"fmt"

	"github.com/Yarik7610/library-backend/user-service/internal/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/model"
	"github.com/Yarik7610/library-backend/user-service/internal/repository"
	"github.com/Yarik7610/library-backend/user-service/internal/utils"
	"gorm.io/gorm"
)

type UserService interface {
	SignUp(user *dto.SignUpUser) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) SignUp(user *dto.SignUpUser) (*model.User, error) {
	foundUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if foundUser != nil {
		return nil, fmt.Errorf("user with such email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password error: %v", err)
	}

	newUser := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	if err = s.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("user create error: %v", err)
	}

	return newUser, nil
}
