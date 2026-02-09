package domain

type User struct {
	ID          uint
	Name        string
	Email       string
	RawPassword string
	IsAdmin     bool
}
