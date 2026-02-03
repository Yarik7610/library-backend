package model

import (
	"time"
)

type Page struct {
	ID        uint `gorm:"primarykey"`
	BookID    uint `gorm:"uniqueIndex:book_id_number_index"`
	Number    uint `gorm:"uniqueIndex:book_id_number_index"`
	Content   string
	CreatedAt time.Time
}
