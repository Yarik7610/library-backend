package postgres

import (
	"context"
	"time"

	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres/model"

	postgresInfrastructure "github.com/Yarik7610/library-backend/user-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, userID uint) (*model.User, error)
	FindByEmail(cxt context.Context, email string) (*model.User, error)
	GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error)
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	name    string
	timeout time.Duration
	db      *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{name: "User(s)", timeout: 500 * time.Millisecond, db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, userID uint) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &user, nil
}

func (r *userRepository) GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var emails []string
	if err := r.db.WithContext(ctx).Model(&model.User{}).Order("email ASC").Where("id IN ?", userIDs).Pluck("email", &emails).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return emails, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, postgresInfrastructure.NewError(err, r.name)
	}
	return count, nil
}
