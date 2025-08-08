package model

import "gorm.io/gorm"

type BookCategory struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);unique;not null"`
}
