package model

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Fullname string
	Books    []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
