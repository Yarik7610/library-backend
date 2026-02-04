package redis

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis/model"
)

func BookWithAuthorModelsToDomains(bookModels []model.BookWithAuthor) []domain.Book {
	bookDomains := make([]domain.Book, len(bookModels))
	for i := range bookModels {
		bookDomains[i] = BookWithAuthorModelToDomain(&bookModels[i])
	}
	return bookDomains
}

func BookWithAuthorModelToDomain(bookModel *model.BookWithAuthor) domain.Book {
	return domain.Book{
		ID:       bookModel.ID,
		Author:   domain.Author(bookModel.Author),
		Title:    bookModel.Title,
		Year:     bookModel.Year,
		Category: bookModel.Category,
	}
}

func BookDomainsToBookWithAuthorModels(bookDomains []domain.Book) []model.BookWithAuthor {
	bookModels := make([]model.BookWithAuthor, len(bookDomains))
	for i := range bookDomains {
		bookModels[i] = BookDomainToBookWithAuthorModel(&bookDomains[i])
	}
	return bookModels
}

func BookDomainToBookWithAuthorModel(bookDomain *domain.Book) model.BookWithAuthor {
	return model.BookWithAuthor{
		ID:       bookDomain.ID,
		Author:   model.Author(bookDomain.Author),
		Title:    bookDomain.Title,
		Year:     bookDomain.Year,
		Category: bookDomain.Category,
	}
}
