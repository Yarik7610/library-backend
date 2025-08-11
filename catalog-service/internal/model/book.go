package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title    string
	Author   string
	Year     int
	Category string `gorm:"unique"`
	Pages    []Page `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
