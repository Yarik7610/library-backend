package mapper

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/dto"
)

func PageDomainToDTO(pageDomain *domain.Page) dto.Page {
	return dto.Page{
		ID:      pageDomain.ID,
		Number:  pageDomain.Number,
		Content: pageDomain.Content,
	}
}

func CreatePageRequestDTOsToDomains(createPageRequestDTOs []dto.CreatePageRequest) []domain.Page {
	pageDomains := make([]domain.Page, len(createPageRequestDTOs))
	for i := range createPageRequestDTOs {
		pageDomains[i] = CreatePageRequestDTOToDomain(&createPageRequestDTOs[i])
	}
	return pageDomains
}

func CreatePageRequestDTOToDomain(createPageRequestDTO *dto.CreatePageRequest) domain.Page {
	return domain.Page{
		Number:  createPageRequestDTO.Number,
		Content: createPageRequestDTO.Content,
	}
}
