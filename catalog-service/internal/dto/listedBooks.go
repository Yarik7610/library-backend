package dto

import "time"

type ListedBooks struct {
	AuthorID uint         `json:"author_id"`
	Fullname string       `json:"fullname"`
	Books    []ListedBook `json:"books"`
}

type ListedBook struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	Category  string    `json:"category"`
}
