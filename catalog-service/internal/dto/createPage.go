package dto

type CreatePage struct {
	Number  int    `json:"number" binding:"required"`
	Content string `json:"content" binding:"required"`
}
