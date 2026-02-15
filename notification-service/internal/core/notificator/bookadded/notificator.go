package bookadded

import (
	"context"
	"encoding/json"

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
	ctx                context.Context
	cancel             context.CancelFunc
}

func NewNotificator(
	logger *logging.Logger,
	bookAddedReader *kafka.Reader,
	emailSender email.Sender,
	subscriptionClient subscription.Client,
) Notificator {
	ctx, cancel := context.WithCancel(context.Background())

	return &notificator{
		logger:             logger,
		bookAddedReader:    bookAddedReader,
		emailSender:        emailSender,
		subscriptionClient: subscriptionClient,
		ctx:                ctx,
		cancel:             cancel,
	}
}

func (n *notificator) Stop() {
	n.cancel()
	n.bookAddedReader.Close()
}

func (n *notificator) Run() {
	const (
		WORKERS_COUNT = 20
		JOBS_MAX_SIZE = 100
	)

	workerPool := n.startWorkerPool(WORKERS_COUNT, JOBS_MAX_SIZE)
	defer workerPool.Stop()

	for {
		select {
		case <-n.ctx.Done():
			n.logger.Info("Stopping book added notificator")
			return
		default:
		}

		message, err := n.bookAddedReader.FetchMessage(n.ctx)
		if err != nil {
			if n.ctx.Err() != nil {
				return
			}
			n.logger.Error("Added book message fetch error", logging.Error(err))
			continue
		}

		if err := n.processMessage(workerPool, message); err != nil {
			if n.ctx.Err() != nil {
				return
			}
			n.logger.Error("Added book message process error", logging.Error(err))
			continue
		}

		if err := n.bookAddedReader.CommitMessages(n.ctx, message); err != nil {
			if n.ctx.Err() != nil {
				return
			}
			n.logger.Error("Added book message commit error", logging.Any("message", message), logging.Error(err))
		}
	}
}

func (n *notificator) startWorkerPool(workersCount, jobsMaxSize int) workerpool.Pool[job.BookAdded] {
	workerPool := workerpool.New[job.BookAdded](workersCount, jobsMaxSize)

	workerPool.Run(func(job job.BookAdded) {
		body := template.ParseBookAddedTemplate(job.AddedBook)

		if err := n.emailSender.Send(body, []string{job.Email}); err != nil {
			n.logger.Error("Added book send mail error", logging.String("email", job.Email), logging.Error(err))
		}
	})

	return workerPool
}

func (n *notificator) processMessage(workerPool workerpool.Pool[job.BookAdded], message kafka.Message) error {
	addedBook, err := n.parseEvent(message.Value)
	if err != nil {
		return err
	}

	emails, err := n.subscriptionClient.GetBookCategorySubscribedUserEmails(n.ctx, addedBook.Category)
	if err != nil {
		return err
	}

	jobs := n.buildJobs(addedBook, emails)
	workerPool.Feed(jobs)
	return nil
}

func (n *notificator) parseEvent(data []byte) (*event.BookAdded, error) {
	var addedBook event.BookAdded
	if err := json.Unmarshal(data, &addedBook); err != nil {
		return nil, err
	}
	return &addedBook, nil
}

func (n *notificator) buildJobs(book *event.BookAdded, emails []string) []job.BookAdded {
	jobs := make([]job.BookAdded, 0, len(emails))
	for _, email := range emails {
		jobs = append(jobs, job.BookAdded{Email: email, AddedBook: book})
	}
	return jobs
}
