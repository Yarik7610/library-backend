package dto

import (
	"encoding/json"
	"time"
)

type AuthorBooksRaw struct {
	AuthorID uint            `json:"author_id"`
	Fullname string          `json:"fullname"`
	Books    json.RawMessage `json:"books"`
}

type AuthorBookRaw struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	Category  string    `json:"category"`
}
