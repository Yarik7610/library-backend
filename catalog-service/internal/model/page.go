package model

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	BookID  uint `json:"book_id" gorm:"uniqueIndex:book_id_number_index"`
	Number  int  `gorm:"uniqueIndex:book_id_number_index"`
	Content string
}
