package pool

import (
	"sync"
)

type Task func() (interface{}, error)

type Pool struct {
	wg      *sync.WaitGroup
	Tasks   chan Task
	Workers chan chan Task
	Results chan interface{}
}

func New(workers, tasks int) *Pool {
	var wg sync.WaitGroup

	pool := &Pool{
		wg:      &wg,
		Tasks:   make(chan Task, tasks),
		Workers: make(chan chan Task, workers),
		Results: make(chan interface{}),
	}

	for i := 0; i < workers; i++ {
		go func() {
			worker := make(chan Task)
			pool.Workers <- worker
			for {
				select {
				case t := <-worker:
					res, err := t()
					if err != nil {
						// what to do if err?
					}
					go func() {
						pool.Results <- res
					}()
					pool.Workers <- worker
					pool.wg.Done()
				}
			}
		}()
	}

	return pool
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Add(t Task) {
	p.Tasks <- t
	p.wg.Add(1)
}

func (p *Pool) Start() {
	go func() {
		for {
			t := <-p.Tasks
			w := <-p.Workers
			w <- t
		}
	}()
}
