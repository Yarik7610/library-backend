package model

import "time"

type Page struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	BookID    uint `json:"book_id" gorm:"uniqueIndex:book_id_number_index"`
	Number    int  `gorm:"uniqueIndex:book_id_number_index"`
	Content   string
}
