package workerpool

import (
	"sync"
)

type Pool[J any] interface {
	Run(fn func(J))
	Feed(jobs []J)
	Stop()
}

type pool[J any] struct {
	workersCount int
	fn           func(J)
	jobs         chan J
	wg           sync.WaitGroup
}

func New[J any](workersCount, jobsMaxSize int) Pool[J] {
	return &pool[J]{
		workersCount: workersCount,
		jobs:         make(chan J, jobsMaxSize),
	}
}

func (p *pool[J]) Run(fn func(J)) {
	p.fn = fn
	for range p.workersCount {
		go p.worker()
	}
}

func (p *pool[J]) Feed(jobs []J) {
	for _, j := range jobs {
		p.wg.Add(1)
		p.jobs <- j
	}
}

func (p *pool[J]) Stop() {
	close(p.jobs)
	p.wg.Wait()
}

func (p *pool[J]) worker() {
	for j := range p.jobs {
		p.fn(j)
		p.wg.Done()
	}
}
