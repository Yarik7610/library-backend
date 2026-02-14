package job

import "github.com/Yarik7610/library-backend-common/broker/kafka/event"

type BookAdded struct {
	AddedBook *event.BookAdded
	Email     string
}
