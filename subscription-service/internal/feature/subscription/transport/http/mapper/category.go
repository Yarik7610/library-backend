package mapper

import (
	"github.com/Yarik7610/library-backend/subscription-service/internal/domain"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http/dto"
)

func UserBookCategoryDomainToDTO(userBookCategoryDomain *domain.UserBookCategory) dto.UserBookCategory {
	return dto.UserBookCategory{
		ID:           userBookCategoryDomain.ID,
		UserID:       userBookCategoryDomain.UserID,
		BookCategory: userBookCategoryDomain.BookCategory,
	}
}

type SubscribeToBookCategoryRequest struct {
	Category string `json:"category" binding:"required"`
}
