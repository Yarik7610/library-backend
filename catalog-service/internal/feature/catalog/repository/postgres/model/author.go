package model

import "time"

type Author struct {
	ID        uint `gorm:"primarykey"`
	Fullname  string
	CreatedAt time.Time
	Books     []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
