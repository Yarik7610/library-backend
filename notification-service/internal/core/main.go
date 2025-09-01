package core

import (
	"context"
	"encoding/json"

	"github.com/Yarik7610/library-backend-common/broker/event"
	"github.com/Yarik7610/library-backend/notification-service/internal/email"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Notificator interface {
	Run()
}

type notificator struct {
	bookAddedReader *kafka.Reader
	emailSender     email.Sender
}

func NewNotificator(bookAddedReader *kafka.Reader, emailSender email.Sender) Notificator {
	return &notificator{
		bookAddedReader: bookAddedReader,
		emailSender:     emailSender,
	}
}

func (n *notificator) Run() {
	workerPool := NewWorkerPool(WORKER_POOL_SIZE, n.emailSender)
	defer workerPool.Stop()

	for {
		m, err := n.bookAddedReader.ReadMessage(context.Background())
		if err != nil {
			zap.S().Errorf("Book added topic read message error: %v", err)
			continue
		}

		var addedBook event.BookAdded
		if err := json.Unmarshal(m.Value, &addedBook); err != nil {
			zap.S().Errorf("Added book unmarshal error: %v", err)
			continue
		}

		emails, err := n.getCategorySubscribersEmails(addedBook.Category)
		if err != nil {
			zap.S().Errorf("Get book's category subscriber's emails error: %v", err)
			continue
		}

		workerPool.Feed(&addedBook, emails)
	}
}
