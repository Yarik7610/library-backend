package query

type ListBooksByCategory struct {
	Page  int    `form:"page" binding:"required,min=1"`
	Count int    `form:"count" binding:"required,min=1,max=100"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}
