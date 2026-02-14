package bookadded

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Yarik7610/library-backend-common/broker/kafka/event"
	"github.com/Yarik7610/library-backend/notification-service/internal/core/job"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/email"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/email/template"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/transport/http/microservice/subscription"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/workerpool"
	"github.com/segmentio/kafka-go"
)

type Notificator interface {
	Run()
	Stop()
}

type notificator struct {
	logger             *logging.Logger
	bookAddedReader    *kafka.Reader
	emailSender        email.Sender
	subscriptionClient subscription.Client
	stop               chan struct{}
}

func NewNotificator(
	logger *logging.Logger,
	bookAddedReader *kafka.Reader,
	emailSender email.Sender,
	subscriptionClient subscription.Client,
) Notificator {
	return &notificator{
		logger:             logger,
		bookAddedReader:    bookAddedReader,
		emailSender:        emailSender,
		subscriptionClient: subscriptionClient,
		stop:               make(chan struct{}),
	}
}

func (n *notificator) Stop() {
	close(n.stop)
}

func (n *notificator) Run() {
	defer n.bookAddedReader.Close()

	const WORKER_POOL_SIZE = 100

	workerPool := n.startWorkerPool(WORKER_POOL_SIZE)
	defer workerPool.Stop()

	for {
		select {
		case <-n.stop:
			n.logger.Info("Stopping book added notificator")
			return
		default:
		}

		message, err := n.readMessage()
		if err != nil {
			continue
		}

		n.processMessage(workerPool, message)
	}
}

func (n *notificator) startWorkerPool(size int) workerpool.Pool[job.BookAdded] {
	workerPool := workerpool.New[job.BookAdded](size)

	workerPool.Run(func(job job.BookAdded) {
		body := template.GetBookAddedEmailTemplate(job.AddedBook)

		if err := n.emailSender.Send(body, []string{job.Email}); err != nil {
			n.logger.Error("Added book send mail error", logging.String("email", job.Email), logging.Error(err))
		}
	})

	return workerPool
}

func (n *notificator) readMessage() (kafka.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	message, err := n.bookAddedReader.ReadMessage(ctx)
	if err != nil {
		n.logger.Error("Added book topic read message error", logging.Error(err))
		return kafka.Message{}, err
	}

	return message, nil
}

func (n *notificator) processMessage(workerPool workerpool.Pool[job.BookAdded], message kafka.Message) {
	addedBook, err := n.parseEvent(message.Value)
	if err != nil {
		return
	}

	emails, err := n.fetchSubscribedUserEmails(addedBook.Category)
	if err != nil {
		return
	}

	jobs := n.buildJobs(addedBook, emails)
	workerPool.Feed(jobs)
}

func (n *notificator) parseEvent(data []byte) (*event.BookAdded, error) {
	var addedBook event.BookAdded
	if err := json.Unmarshal(data, &addedBook); err != nil {
		n.logger.Error("Added book unmarshal error", logging.Error(err))
		return nil, err
	}
	return &addedBook, nil
}

func (n *notificator) fetchSubscribedUserEmails(bookCategory string) ([]string, error) {
	emails, err := n.subscriptionClient.GetBookCategorySubscribedUserEmails(context.Background(), bookCategory)
	if err != nil {
		n.logger.Error("Get book category subscribed user emails error", logging.Error(err))
		return nil, err
	}
	return emails, nil
}

func (n *notificator) buildJobs(book *event.BookAdded, emails []string) []job.BookAdded {
	jobs := make([]job.BookAdded, 0, len(emails))
	for _, email := range emails {
		jobs = append(jobs, job.BookAdded{Email: email, AddedBook: book})
	}
	return jobs
}
