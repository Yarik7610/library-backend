package model

import (
	"time"
)

type Book struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	AuthorID  uint      `json:"author_id" gorm:"uniqueIndex:author_id_title_index"`
	Title     string    `json:"title" gorm:"uniqueIndex:author_id_title_index"`
	Year      int       `json:"year"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	Pages     []Page    `json:"pages,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
