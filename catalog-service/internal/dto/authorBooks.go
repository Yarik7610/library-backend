package dto

type AuthorBooks struct {
	AuthorID uint            `json:"author_id"`
	Fullname string          `json:"fullname"`
	Books    []AuthorBookRaw `json:"books"`
}
