package event

type BookAdded struct {
	ID       uint   `json:"id"`
	AuthorID uint   `json:"author_id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Category string `json:"category"`
}
