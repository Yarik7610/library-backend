package model

import "time"

type UserCategory struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"index;uniqueIndex:user_id_category_index"`
	Category  string    `json:"category" gorm:"index;uniqueIndex:user_id_category_index"`
	CreatedAt time.Time `json:"created_at"`
}
