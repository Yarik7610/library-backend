package dto

type CreatePageRequest struct {
	Number  uint   `json:"number" binding:"required"`
	Content string `json:"content" binding:"required"`
}
