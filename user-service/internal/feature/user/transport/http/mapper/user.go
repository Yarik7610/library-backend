package mapper

import (
	"github.com/Yarik7610/library-backend/user-service/internal/domain"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http/dto"
)

func UserDomainToDTO(userDomain *domain.User) dto.User {
	return dto.User{
		ID:      userDomain.ID,
		Name:    userDomain.Name,
		Email:   userDomain.Email,
		IsAdmin: userDomain.IsAdmin,
	}
}

func SignUpUserRequestDTOToDomain(signUpUserRequestDTO *dto.SignUpUserRequest) domain.User {
	return domain.User{
		Name:        signUpUserRequestDTO.Name,
		Email:       signUpUserRequestDTO.Email,
		RawPassword: signUpUserRequestDTO.Password,
	}
}

func SignInUserRequestDTOToDomain(signInUserRequestDTO *dto.SignInUserRequest) domain.User {
	return domain.User{
		Email:       signInUserRequestDTO.Email,
		RawPassword: signInUserRequestDTO.Password,
	}
}
