package dto

type Page struct {
	ID      uint   `json:"id"`
	Number  uint   `json:"number"`
	Content string `json:"content"`
}

type CreatePageRequest struct {
	Number  uint   `json:"number" binding:"required"`
	Content string `json:"content" binding:"required"`
}
