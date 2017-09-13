// package pool provides a simple worker pool
package pool

import (
	"sync"
)

// Task describes a function that can return any value, followed by an error
type Task func() (interface{}, error)

// Pool holds a collection of worker goroutines,
// executing all of the incoming Tasks as workers become available
type Pool struct {
	wg      *sync.WaitGroup
	Tasks   chan Task
	Workers chan chan Task
	Results chan interface{}
}

// New returns a new Pool
//
// workers sets the concurrency at which to work
//
// tasks sets the size of our Task queue
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

// Wait for all Tasks to be processed
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Add a new Task to the Pool
func (p *Pool) Add(t Task) {
	p.Tasks <- t
	p.wg.Add(1)
}

// Start working through the Task queue
func (p *Pool) Start() {
	go func() {
		for {
			t := <-p.Tasks
			w := <-p.Workers
			w <- t
		}
	}()
}
