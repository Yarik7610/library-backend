package model

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint `gorm:"primarykey;uniqueIndex:book_id_number_index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string
	Author    string
	Year      int
	Category  string `gorm:"uniqueIndex:book_id_number_index"`
	Pages     []Page `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
