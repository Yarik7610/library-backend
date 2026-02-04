package redis

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis/model"
)

func BookModelsToDomains(bookModels []model.Book) []domain.Book {
	bookDomains := make([]domain.Book, len(bookModels))
	for i := range bookModels {
		bookDomains[i] = BookModelToDomain(&bookModels[i])
	}
	return bookDomains
}

func BookModelToDomain(bookModel *model.Book) domain.Book {
	return domain.Book{
		ID:       bookModel.ID,
		Author:   domain.Author(bookModel.Author),
		Title:    bookModel.Title,
		Year:     bookModel.Year,
		Category: bookModel.Category,
	}
}

func BookDomainsToModels(bookDomains []domain.Book) []model.Book {
	bookModels := make([]model.Book, len(bookDomains))
	for i := range bookDomains {
		bookModels[i] = BookDomainToModel(&bookDomains[i])
	}
	return bookModels
}

func BookDomainToModel(bookDomain *domain.Book) model.Book {
	return model.Book{
		ID:       bookDomain.ID,
		Author:   model.Author(bookDomain.Author),
		Title:    bookDomain.Title,
		Year:     bookDomain.Year,
		Category: bookDomain.Category,
	}
}
