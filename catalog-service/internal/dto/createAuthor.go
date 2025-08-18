package dto

type CreateAuthor struct {
	Fullname string `json:"fullname" binding:"required"`
}
