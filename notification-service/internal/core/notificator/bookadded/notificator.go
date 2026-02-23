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
	Run(ctx context.Context)
	Stop(ctx context.Context)
}

type notificator struct {
	logger             *logging.Logger
	bookAddedReader    *kafka.Reader
	emailSender        email.Sender
	subscriptionClient subscription.Client
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
	}
}

func (n *notificator) Stop(ctx context.Context) {
	n.bookAddedReader.Close()
}

func (n *notificator) Run(ctx context.Context) {
	const (
		WORKERS_COUNT = 20
		JOBS_MAX_SIZE = 100
	)

	workerPool := n.startWorkerPool(ctx, WORKERS_COUNT, JOBS_MAX_SIZE)
	defer workerPool.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		message, err := n.bookAddedReader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			n.logger.Error(ctx, "Added book message fetch error", logging.Error(err))
			continue
		}

		if err := n.processMessage(ctx, workerPool, message); err != nil {
			if ctx.Err() != nil {
				return
			}
			n.logger.Error(ctx, "Added book message process error", logging.Error(err))
			continue
		}

		if err := n.bookAddedReader.CommitMessages(ctx, message); err != nil {
			if ctx.Err() != nil {
				return
			}
			n.logger.Error(ctx, "Added book message commit error", logging.Any("message", message), logging.Error(err))
		}
	}
}

func (n *notificator) startWorkerPool(ctx context.Context, workersCount, jobsMaxSize int) workerpool.Pool[job.BookAdded] {
	workerPool := workerpool.New[job.BookAdded](workersCount, jobsMaxSize)

	workerPool.Run(func(job job.BookAdded) {
		body := template.ParseBookAddedTemplate(job.AddedBook)

		if err := n.emailSender.Send(body, []string{job.Email}); err != nil {
			n.logger.Error(ctx, "Added book send mail error", logging.String("email", job.Email), logging.Error(err))
		}
	})

	return workerPool
}

func (n *notificator) processMessage(ctx context.Context, workerPool workerpool.Pool[job.BookAdded], message kafka.Message) error {
	addedBook, err := n.parseEvent(message.Value)
	if err != nil {
		return err
	}

	emails, err := n.subscriptionClient.GetBookCategorySubscribedUserEmails(ctx, addedBook.Category)
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
