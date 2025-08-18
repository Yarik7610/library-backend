package model

import "time"

type Author struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Fullname  string
	Books     []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
