package query

type GetBookPage struct {
	PageNumber int `form:"page" binding:"required,min=1"`
}
