package model

type UserBookCategory struct {
	UserID     uint `gorm:"primaryKey"`
	CategoryID uint `gorm:"primaryKey"`
}
