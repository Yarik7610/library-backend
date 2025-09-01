package core

import (
	"sync"

	"github.com/Yarik7610/library-backend-common/broker/event"
)

type WorkerPool[I any] interface {
	Run(fn func(I))
	Stop()
	Feed(jobs []I)
}

type emailJob struct {
	addedBook *event.BookAdded
	email     string
}

type workerPool[I any] struct {
	size int
	fn   func(I)
	wg   sync.WaitGroup
	jobs chan I
}

func NewWorkerPool[I any](size int) WorkerPool[I] {
	pool := &workerPool[I]{
		size: size,
		jobs: make(chan I, size),
	}

	return pool
}

func (p *workerPool[I]) Run(fn func(I)) {
	p.fn = fn
	for range p.size {
		go p.worker()
	}
}

func (p *workerPool[I]) Stop() {
	close(p.jobs)
}

func (p *workerPool[I]) Feed(tasks []I) {
	for _, task := range tasks {
		p.wg.Add(1)
		p.jobs <- task
	}
	p.wg.Wait()
}

func (p *workerPool[I]) worker() {
	for job := range p.jobs {
		p.fn(job)
		p.wg.Done()
	}
}
