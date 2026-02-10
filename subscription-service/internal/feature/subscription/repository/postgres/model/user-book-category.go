package model

import "time"

type UserBookCategory struct {
	ID           uint   `gorm:"primarykey"`
	UserID       uint   `gorm:"index;uniqueIndex:user_id_book_category_index"`
	BookCategory string `gorm:"index;uniqueIndex:user_id_book_category_index"`
	CreatedAt    time.Time
}
