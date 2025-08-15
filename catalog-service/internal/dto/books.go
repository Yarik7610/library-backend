package dto

import "time"

type Books struct {
	AuthorID uint   `json:"author_id"`
	Fullname string `json:"fullname"`
	Books    []Book `json:"books"`
}

type Book struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	Category  string    `json:"category"`
}
