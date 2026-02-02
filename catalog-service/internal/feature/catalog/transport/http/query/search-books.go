package query

type SearchBooks struct {
	Author string `form:"author"`
	Title  string `form:"title"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
	Page   uint   `form:"page" binding:"required,min=1"`
	Count  uint   `form:"count" binding:"required,min=1,max=100"`
}
