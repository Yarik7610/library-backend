package service

import (
	"fmt"
	"net/http"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/user-service/internal/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/model"
	"github.com/Yarik7610/library-backend/user-service/internal/repository"
	"github.com/Yarik7610/library-backend/user-service/internal/utils"
)

type UserService interface {
	SignUp(user *dto.SignUpUser) (*model.User, *custom.Err)
	SignIn(user *dto.SignInUser) (string, *custom.Err)
	Me(userID uint) (*model.User, *custom.Err)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) SignUp(user *dto.SignUpUser) (*model.User, *custom.Err) {
	foundUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	if foundUser != nil {
		return nil, custom.NewErr(http.StatusBadRequest, "user with such email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, custom.NewErr(http.StatusBadRequest, fmt.Sprintf("hash password error: %v", err))
	}

	newUser := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	if err = s.userRepo.Create(newUser); err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, fmt.Sprintf("user create error: %v", err))
	}

	return newUser, nil
}

func (s *userService) SignIn(user *dto.SignInUser) (string, *custom.Err) {
	foundUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return "", custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	if foundUser == nil {
		return "", custom.NewErr(http.StatusBadRequest, "wrong email or password")
	}

	if !utils.CompareHashAndPasword(foundUser.Password, user.Password) {
		return "", custom.NewErr(http.StatusBadRequest, "wrong email or password")
	}

	token, err := utils.CreateJWTToken(foundUser.ID, foundUser.IsAdmin)
	if err != nil {
		return "", custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	return token, nil
}

func (s *userService) Me(userID uint) (*model.User, *custom.Err) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}
