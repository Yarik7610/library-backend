package postgres

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
)

func BookModelToDomain(bookModel *model.Book) domain.Book {
	return domain.Book{
		ID:       bookModel.ID,
		Title:    bookModel.Title,
		Year:     bookModel.Year,
		Category: bookModel.Category,
	}
}

func BookWithAuthorModelToDomain(bookWithAuthorModel *model.BookWithAuthor) domain.Book {
	return domain.Book{
		ID: bookWithAuthorModel.ID,
		Author: domain.Author{
			ID:       bookWithAuthorModel.AuthorID,
			Fullname: bookWithAuthorModel.AuthorFullname,
		},
		Title:    bookWithAuthorModel.Title,
		Year:     bookWithAuthorModel.Year,
		Category: bookWithAuthorModel.Category,
	}
}

func BookWithAuthorModelsToDomains(bookWithAuthorModels []model.BookWithAuthor) []domain.Book {
	bookDomains := make([]domain.Book, len(bookWithAuthorModels))
	for i := range bookWithAuthorModels {
		bookDomains[i] = BookWithAuthorModelToDomain(&bookWithAuthorModels[i])
	}
	return bookDomains
}
