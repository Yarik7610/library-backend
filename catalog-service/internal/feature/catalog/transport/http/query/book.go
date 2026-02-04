package query

type ListBooksByCategory struct {
	Page  uint   `form:"page,default=1" binding:"min=1"`
	Count uint   `form:"count,default=20" binding:"min=1,max=100"`
	Sort  string `form:"sort,default=title"`
	Order string `form:"order,default=asc"`
}

type SearchBooks struct {
	Author string `form:"author"`
	Title  string `form:"title"`
	Page   uint   `form:"page,default=1" binding:"min=1"`
	Count  uint   `form:"count,default=20" binding:"min=1,max=100"`
	Sort   string `form:"sort,default=title"`
	Order  string `form:"order,default=asc"`
}
