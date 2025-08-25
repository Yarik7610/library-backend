package model

import "time"

type UserCategory struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}
