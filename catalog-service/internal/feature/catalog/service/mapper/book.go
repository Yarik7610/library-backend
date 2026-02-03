package mapper

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
