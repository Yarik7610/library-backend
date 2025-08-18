package query

type GetBookPage struct {
	PageNumber uint `form:"page" binding:"required,min=1"`
}
