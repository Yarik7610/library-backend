package core

import (
	"context"
	"encoding/json"

	"github.com/Yarik7610/library-backend-common/broker/event"
	"github.com/Yarik7610/library-backend/notification-service/internal/email"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Controller interface {
	Start()
}

type controller struct {
	bookAddedReader *kafka.Reader
	emailSender     email.Sender
}

func NewController(bookAddedReader *kafka.Reader, emailSender email.Sender) Controller {
	return &controller{
		bookAddedReader: bookAddedReader,
		emailSender:     emailSender,
	}
}

func (c *controller) Start() {
	for {
		m, err := c.bookAddedReader.ReadMessage(context.Background())
		if err != nil {
			zap.S().Errorf("Book added topic read message error: %v", err)
			continue
		}

		var addedBook event.BookAdded
		if err := json.Unmarshal(m.Value, &addedBook); err != nil {
			zap.S().Errorf("Added book unmarshal error: %v", err)
			continue
		}

		emails, err := c.getCategorySubscribersEmails(addedBook.Category)
		if err != nil {
			zap.S().Errorf("Get book's category subscribers emails error: %v", err)
			continue
		}

		// body := fmt.Sprintf("Hello! New books arrival in %q category", utils.Capitalize("horror"))
		zap.S().Debug("EMAILS", emails)

		// c.emailSender.Send(body, []string{config.Data.Mail})
	}
}
