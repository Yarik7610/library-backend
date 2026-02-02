package dto

type CreateAuthorRequest struct {
	Fullname string `json:"fullname" binding:"required"`
}
