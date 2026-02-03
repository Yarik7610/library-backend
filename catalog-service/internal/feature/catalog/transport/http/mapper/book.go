package mapper

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/dto"
)

func BookDomainsToDTOs(bookDomains []domain.Book) []dto.Book {
	bookDTOs := make([]dto.Book, len(bookDomains))
	for i := range bookDomains {
		bookDTOs[i] = BookDomainToDTO(&bookDomains[i])
	}
	return bookDTOs
}

func BookDomainToDTO(bookDomain *domain.Book) dto.Book {
	return dto.Book{
		ID:       bookDomain.ID,
		Author:   dto.Author(bookDomain.Author),
		Title:    bookDomain.Title,
		Year:     bookDomain.Year,
		Category: bookDomain.Category,
	}
}

func AddBookRequestToDomain(addBookRequestDTO *dto.AddBookRequest) domain.Book {
	return domain.Book{
		Author:   domain.Author{ID: addBookRequestDTO.AuthorID},
		Title:    addBookRequestDTO.Title,
		Year:     addBookRequestDTO.Year,
		Category: addBookRequestDTO.Category,
		Pages:    CreatePageRequestDTOsToDomains(addBookRequestDTO.Pages),
	}
}
