package mapper

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/dto"
)

func AuthorDomainToDTO(authorDomain *domain.Author) dto.Author {
	return dto.Author{
		ID:       authorDomain.ID,
		Fullname: authorDomain.Fullname,
	}
}

func CreateAuthorRequestDTOToDomain(createAuthorRequestDTO *dto.CreateAuthorRequest) domain.Author {
	return domain.Author{Fullname: createAuthorRequestDTO.Fullname}
}
