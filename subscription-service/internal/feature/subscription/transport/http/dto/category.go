package dto

type Create struct {
	Category string `json:"category" binding:"required"`
}
