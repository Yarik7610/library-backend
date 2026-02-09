package model

import "time"

type User struct {
	ID             uint `gorm:"primarykey"`
	Name           string
	Email          string `gorm:"unique"`
	HashedPassword string `gorm:"column:password"`
	IsAdmin        bool
	CreatedAt      time.Time
}
