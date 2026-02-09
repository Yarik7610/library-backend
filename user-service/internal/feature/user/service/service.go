package service

import (
	"context"

	"github.com/Yarik7610/library-backend/user-service/internal/domain"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres"
	mapper "github.com/Yarik7610/library-backend/user-service/internal/feature/user/service/mapper/postgres"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/errs"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/jwt"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/password"
)

type UserService interface {
	SignUp(ctx context.Context, userDomain *domain.User) error
	SignIn(ctx context.Context, userDomain *domain.User) (*domain.Token, error)
	GetMe(ctx context.Context, userID uint) (*domain.User, error)
	GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error)
}

type userService struct {
	userRepository postgres.UserRepository
}

func NewUserService(userRepository postgres.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) SignUp(ctx context.Context, userDomain *domain.User) error {
	userModel, err := mapper.UserDomainToModel(userDomain)
	if err != nil {
		return err
	}

	if err = s.userRepository.Create(ctx, &userModel); err != nil {
		return err
	}

	userDomain.ID = userModel.ID
	userDomain.IsAdmin = userModel.IsAdmin
	return nil
}

func (s *userService) SignIn(ctx context.Context, userDomain *domain.User) (*domain.Token, error) {
	foundUser, err := s.userRepository.FindByEmail(ctx, userDomain.Email)
	if err != nil {
		return nil, err
	}

	if !password.CompareHashAndRaw(foundUser.HashedPassword, userDomain.RawPassword) {
		return nil, errs.NewBadRequestError("Wrong email or password")
	}

	accessToken, err := jwt.Create(foundUser.ID, foundUser.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &domain.Token{AccessToken: accessToken}, nil
}

func (s *userService) GetMe(ctx context.Context, userID uint) (*domain.User, error) {
	userModel, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userDomain := mapper.UserModelToDomain(userModel)
	return &userDomain, nil
}

func (s *userService) GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error) {
	emails, err := s.userRepository.GetEmailsByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return emails, nil
}
