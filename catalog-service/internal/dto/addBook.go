package dto

type AddBook struct {
	AuthorID uint         `json:"author_id" binding:"required,min=1"`
	Title    string       `json:"password" binding:"required"`
	Year     int          `json:"year" binding:"required"`
	Category string       `json:"category" binding:"required"`
	Pages    []CreatePage `json:"pages"`
}
