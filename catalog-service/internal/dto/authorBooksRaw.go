package dto

import (
	"encoding/json"
)

type AuthorBooksRaw struct {
	AuthorID uint            `json:"author_id"`
	Fullname string          `json:"fullname"`
	Books    json.RawMessage `json:"books"`
}
