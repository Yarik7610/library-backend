package dto

type AddBookRequest struct {
	AuthorID uint                `json:"author_id" binding:"required,min=1"`
	Title    string              `json:"title" binding:"required"`
	Year     int                 `json:"year" binding:"required"`
	Category string              `json:"category" binding:"required"`
	Pages    []CreatePageRequest `json:"pages"`
}
