package seed

import (
	"context"

	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/password"

	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres/model"
	"go.uber.org/zap"
)

func Admin(userRepository postgres.UserRepository) {
	ctx := context.Background()

	usersCount, err := userRepository.Count(ctx)
	if err != nil {
		zap.S().Fatalf("Failed to count users for seed need: %v", err)
	}

	if usersCount != 0 {
		zap.S().Info("Seeded admin already exist, skip seeding")
		return
	}

	zap.S().Info("No admin found, start seeding...")
	if err := seedAdmin(ctx, userRepository); err != nil {
		zap.S().Fatalf("Admin seed error: %v", err)
	}
	zap.S().Info("Successfully seeded admin")
}

func seedAdmin(ctx context.Context, userRepository postgres.UserRepository) error {
	hashedPassword, err := password.GenerateHash("admin")
	if err != nil {
		return err
	}

	admin := model.User{
		Name:     "admin",
		Email:    config.Data.Mail,
		Password: hashedPassword,
		IsAdmin:  true,
	}
	return userRepository.Create(ctx, &admin)
}
