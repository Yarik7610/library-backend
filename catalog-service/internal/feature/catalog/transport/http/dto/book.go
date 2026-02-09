package dto

type Book struct {
	ID       uint   `json:"id"`
	Author   Author `json:"author"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Category string `json:"category"`
}

type AddBookRequest struct {
	AuthorID uint                `json:"authorId" binding:"required,min=1"`
	Title    string              `json:"title" binding:"required"`
	Year     int                 `json:"year" binding:"required"`
	Category string              `json:"category" binding:"required"`
	Pages    []CreatePageRequest `json:"pages"`
}

type BookViews struct {
	Views int64 `json:"views"`
}
