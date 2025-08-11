package model

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	BookID  uint
	Number  int `gorm:"unique"`
	Content string
}
