package model

import "time"

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	IsAdmin   bool
}
