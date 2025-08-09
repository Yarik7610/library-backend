package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
	IsAdmin  bool

	SubscribedBookCategories []*BookCategory `gorm:"many2many:user_book_categories;"`
}
