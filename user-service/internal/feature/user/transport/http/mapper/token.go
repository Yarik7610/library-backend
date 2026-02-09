package mapper

import (
	"github.com/Yarik7610/library-backend/user-service/internal/domain"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http/dto"
)

func TokenDomainToDTO(tokenDomain *domain.Token) dto.Token {
	return dto.Token{AccessToken: tokenDomain.AccessToken}
}
