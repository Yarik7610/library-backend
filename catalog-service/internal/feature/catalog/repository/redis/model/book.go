package model

type Book struct {
	ID       uint
	Author   Author
	Title    string
	Year     int
	Category string
}
