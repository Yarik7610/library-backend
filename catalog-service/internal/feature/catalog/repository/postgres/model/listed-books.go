package model

import "encoding/json"

type ListedBooks struct {
	AuthorID uint
	Fullname string
	Books    json.RawMessage
}
