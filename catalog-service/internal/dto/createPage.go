package dto

type CreatePage struct {
	Number  uint   `json:"number" binding:"required"`
	Content string `json:"content" binding:"required"`
}
