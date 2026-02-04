package model

type BookWithAuthor struct {
	ID       uint
	Author   Author
	Title    string
	Year     int
	Category string
}
