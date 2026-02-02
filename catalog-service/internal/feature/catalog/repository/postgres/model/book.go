package model

import (
	"time"
)

type Book struct {
	ID        uint   `gorm:"primarykey"`
	AuthorID  uint   `gorm:"uniqueIndex:author_id_title_index"`
	Title     string `gorm:"uniqueIndex:author_id_title_index"`
	Year      int
	Category  string
	CreatedAt time.Time
	Pages     []Page `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
