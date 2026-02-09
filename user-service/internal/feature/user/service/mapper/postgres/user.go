package mapper

import (
	"github.com/Yarik7610/library-backend/user-service/internal/domain"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/repository/postgres/model"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/password"
)

func UserModelToDomain(userModel *model.User) domain.User {
	return domain.User{
		ID:      userModel.ID,
		Name:    userModel.Email,
		Email:   userModel.Email,
		IsAdmin: userModel.IsAdmin,
	}
}

func UserDomainToModel(userDomain *domain.User) (model.User, error) {
	hashedPassword, err := password.GenerateHash(userDomain.RawPassword)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:             userDomain.ID,
		Name:           userDomain.Name,
		Email:          userDomain.Email,
		IsAdmin:        userDomain.IsAdmin,
		HashedPassword: hashedPassword,
	}, nil
}
