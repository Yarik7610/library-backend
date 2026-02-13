package dto

type UserBookCategory struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"userId"`
	BookCategory string `json:"bookCategory"`
}

type SubscribeToBookCategoryRequest struct {
	BookCategory string `json:"bookCategory" binding:"required"`
}
