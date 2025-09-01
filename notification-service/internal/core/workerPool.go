package core

import (
	"sync"

	"github.com/Yarik7610/library-backend-common/broker/event"
	"github.com/Yarik7610/library-backend/notification-service/internal/email"
	"go.uber.org/zap"
)

type WorkerPool interface {
	Stop()
	Feed(addedBook *event.BookAdded, emails []string)
}

type job struct {
	addedBook *event.BookAdded
	email     string
	wg        *sync.WaitGroup
}

type workerPool struct {
	size        int
	emailSender email.Sender
	jobs        chan job
}

func NewWorkerPool(size int, emailSender email.Sender) WorkerPool {
	pool := &workerPool{
		size:        size,
		emailSender: emailSender,
		jobs:        make(chan job, size),
	}

	pool.run()
	return pool
}

func (p *workerPool) Stop() {
	close(p.jobs)
}

func (p *workerPool) Feed(addedBook *event.BookAdded, emails []string) {
	var wg sync.WaitGroup
	for _, email := range emails {
		wg.Add(1)
		p.jobs <- job{addedBook: addedBook, email: email, wg: &wg}
	}
	wg.Wait()
}

func (p *workerPool) run() {
	for range p.size {
		go p.worker()
	}
}

func (p *workerPool) worker() {
	for job := range p.jobs {
		body := email.FormatBookAddedEmailHTML(job.addedBook)
		err := p.emailSender.Send(body, []string{job.email})
		if err != nil {
			zap.S().Errorf("Mail send to %v error: %v", job.email, err)
		}
		job.wg.Done()
	}
}
