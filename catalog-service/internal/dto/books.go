package dto

type Books struct {
	AuthorID uint      `json:"author_id"`
	Fullname string    `json:"fullname"`
	Books    []BookRaw `json:"books"`
}
