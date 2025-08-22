package model

import "time"

type Page struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	BookID    uint      `json:"book_id" gorm:"uniqueIndex:book_id_number_index"`
	Number    uint      `json:"number" gorm:"uniqueIndex:book_id_number_index"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
