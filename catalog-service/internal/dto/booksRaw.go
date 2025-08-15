package dto

import (
	"encoding/json"
)

type BooksRaw struct {
	AuthorID uint            `json:"author_id"`
	Fullname string          `json:"fullname"`
	Books    json.RawMessage `json:"books"`
}
