package domain

type ListedBooks struct {
	AuthorID uint
	Fullname string
	Books    []ListedBook
}

type ListedBook struct {
	ID       uint
	Title    string
	Year     int
	Category string
}
