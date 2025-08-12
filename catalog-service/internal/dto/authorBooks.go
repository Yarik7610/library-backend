package dto

import "github.com/Yarik7610/library-backend/catalog-service/internal/model"

type AuthorBooks struct {
	AuthorID uint         `json:"author_id"`
	Fullname string       `json:"fullname"`
	Books    []model.Book `json:"books"`
}
