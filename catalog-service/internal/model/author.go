package model

import "time"

type Author struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Fullname  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
	Books     []Book    `json:"books,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
