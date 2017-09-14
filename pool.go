// Package pool provides a simple worker pool
package pool

import (
	"sync"
)

// Task describes a function that returns a value of any type & an error
type Task func() (interface{}, error)

// Result describes the values coming out of the Pool results channel
type Result struct {
	Value interface{}
	Err   error
}

// Pool holds a collection of worker goroutines,
// executing all of the incoming Tasks as workers become available
type Pool struct {
	wg      *sync.WaitGroup
	Tasks   chan Task
	Workers chan chan Task
	Results chan Result
}

// New returns a new Pool.
//
// workers sets the concurrency at which to work.
//
// queue sets the size of our Task queue.
func New(workers, queue int) *Pool {
	var wg sync.WaitGroup

	pool := &Pool{
		wg:      &wg,
		Tasks:   make(chan Task, queue),
		Workers: make(chan chan Task, workers),
		Results: make(chan Result),
	}

	for i := 0; i < workers; i++ {
		go func() {
			worker := make(chan Task)
			pool.Workers <- worker
			for {
				select {
				case t := <-worker:
					res, err := t()
					go func() {
						pool.Results <- Result{Value: res, Err: err}
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
