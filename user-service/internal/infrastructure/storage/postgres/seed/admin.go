package seed

import (
	"context"

	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/password"

	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres/model"
)

func Admin(config *config.Config, userRepository postgres.UserRepository) error {
	ctx := context.Background()

	usersCount, err := userRepository.Count(ctx)
	if err != nil {
		return err
	}

	if usersCount != 0 {
		return nil
	}

	return seedAdmin(ctx, config, userRepository)
}

func seedAdmin(ctx context.Context, config *config.Config, userRepository postgres.UserRepository) error {
	hashedPassword, err := password.GenerateHash("admin")
	if err != nil {
		return err
	}

	admin := model.User{
		Name:           "admin",
		Email:          config.Mail,
		HashedPassword: hashedPassword,
		IsAdmin:        true,
	}
	return userRepository.Create(ctx, &admin)
}
