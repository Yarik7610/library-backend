package dto

type SubscribeCategory struct {
	Category string `json:"category" binding:"required"`
}
