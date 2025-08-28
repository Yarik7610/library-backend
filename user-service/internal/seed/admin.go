package seed

import (
	"github.com/Yarik7610/library-backend/user-service/config"
	"github.com/Yarik7610/library-backend/user-service/internal/model"
	"github.com/Yarik7610/library-backend/user-service/internal/repository"
	"github.com/Yarik7610/library-backend/user-service/internal/utils"
	"go.uber.org/zap"
)

func Admin(userRepository repository.UserRepository) {
	usersCount, err := userRepository.CountUsers()
	if err != nil {
		zap.S().Fatalf("Failed to count users for seed need: %v", err)
	}

	if usersCount != 0 {
		zap.S().Info("Seeded admin already exist, skip seeding")
		return
	}

	zap.S().Info("No admin found, start seeding...")
	if err := seedAdmin(userRepository); err != nil {
		zap.S().Fatalf("Admin seed error: %v", err)
	}
	zap.S().Info("Successfully seeded admin")
}

func seedAdmin(userRepository repository.UserRepository) error {
	hashedPassword, err := utils.HashPassword("admin")
	if err != nil {
		return err
	}

	admin := model.User{
		Name:     "admin",
		Email:    config.Data.Mail,
		Password: hashedPassword,
		IsAdmin:  true,
	}
	return userRepository.Create(&admin)
}
