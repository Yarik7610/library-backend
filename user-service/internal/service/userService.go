package service

import (
	"fmt"
	"net/http"

	"github.com/Yarik7610/library-backend/user-service/config"
	"github.com/Yarik7610/library-backend/user-service/internal/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/model"
	"github.com/Yarik7610/library-backend/user-service/internal/repository"
	"github.com/Yarik7610/library-backend/user-service/internal/utils"
	apperror "github.com/Yarik7610/library-backend/user-service/pkg/app-error"
)

type UserService interface {
	SignUp(user *dto.SignUpUser) (*model.User, *apperror.Err)
	SignIn(user *dto.SignInUser) (string, *apperror.Err)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) SignUp(user *dto.SignUpUser) (*model.User, *apperror.Err) {
	foundUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, err.Error())
	}

	if foundUser != nil {
		return nil, apperror.New(http.StatusBadRequest, "user with such email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, apperror.New(http.StatusBadRequest, fmt.Sprintf("hash password error: %v", err))
	}

	newUser := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	if err = s.userRepo.Create(newUser); err != nil {
		return nil, apperror.New(http.StatusInternalServerError, fmt.Sprintf("user create error: %v", err))
	}

	return newUser, nil
}

func (s *userService) SignIn(user *dto.SignInUser) (string, *apperror.Err) {
	foundUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return "", apperror.New(http.StatusInternalServerError, err.Error())
	}

	if foundUser == nil {
		return "", apperror.New(http.StatusBadRequest, "wrong email or password")
	}

	if !utils.CompareHashAndPasword(foundUser.Password, user.Password) {
		return "", apperror.New(http.StatusBadRequest, "wrong email or password")
	}

	token, err := utils.CreateJWTToken(foundUser.ID, foundUser.IsAdmin, config.Data.JWTSecret)
	if err != nil {
		return "", apperror.New(http.StatusInternalServerError, err.Error())
	}

	return token, nil
}
