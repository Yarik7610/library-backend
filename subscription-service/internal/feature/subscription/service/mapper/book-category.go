package mapper

import (
	"github.com/Yarik7610/library-backend/subscription-service/internal/domain"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres/model"
)

func UserBookCategoryToDomain(userBookCategoryModel *model.UserBookCategory) domain.UserBookCategory {
	return domain.UserBookCategory{
		ID:           userBookCategoryModel.ID,
		UserID:       userBookCategoryModel.UserID,
		BookCategory: userBookCategoryModel.BookCategory,
	}
}
