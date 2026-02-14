package workerpool

import (
	"sync"
)

type Pool[J any] interface {
	Run(fn func(J))
	Stop()
	Feed(jobs []J)
}

type pool[J any] struct {
	size int
	fn   func(J)
	jobs chan J
	wg   sync.WaitGroup
}

func New[J any](size int) Pool[J] {
	return &pool[J]{
		size: size,
		jobs: make(chan J, size),
	}
}

func (p *pool[J]) Run(fn func(J)) {
	p.fn = fn
	for range p.size {
		go p.worker()
	}
}

func (p *pool[J]) Stop() {
	close(p.jobs)
}

func (p *pool[J]) Feed(tasks []J) {
	for _, task := range tasks {
		p.wg.Add(1)
		p.jobs <- task
	}
	p.wg.Wait()
}

func (p *pool[J]) worker() {
	for job := range p.jobs {
		p.fn(job)
		p.wg.Done()
	}
}
