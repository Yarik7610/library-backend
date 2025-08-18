package model

import (
	"time"
)

type Book struct {
	ID        uint `gorm:"primarykey;uniqueIndex:book_id_number_index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	AuthorID  uint `json:"author_id"`
	Year      int
	Category  string `gorm:"uniqueIndex:book_id_number_index"`
	Pages     []Page `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"pages,omitempty"`
}
