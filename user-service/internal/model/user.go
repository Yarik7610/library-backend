package model

import (
	"github.com/Yarik7610/library-backend/catalog-service/pkg/model"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	IsAdmin  bool   `gorm:"default:false"`

	SubscribedBookCategories []model.BookCategory
}
