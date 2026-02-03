package dto

type Author struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
}

type CreateAuthorRequest struct {
	Fullname string `json:"fullname" binding:"required"`
}
