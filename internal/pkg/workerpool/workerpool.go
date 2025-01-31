// Package workerpool provides a simple worker pool implementation for concurrent job processing.
// It allows submitting jobs to a pool of workers, handles synchronization, and supports custom job handlers.
package workerpool

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type jobCh[T any] chan T
type jobHandler[T any] func(T) error

// pool represents the worker pool that manages workers, jobs, and synchronization.
type pool[T any] struct {
	numOfWorkers int            // number of workers in pool
	jobCh        jobCh[T]       // chanel with jobs
	jobHandler   jobHandler[T]  // job handler
	wg           sync.WaitGroup // synchronizer of workers jobs
}

func New[T any](numOfWorkers, jobChSize int, handler jobHandler[T]) *pool[T] {
	// Create a channel to hold jobs
	jobCh := make(chan T, jobChSize)

	// Initialize the pool with the specified parameters
	p := pool[T]{
		jobCh:        jobCh,
		numOfWorkers: numOfWorkers,
		jobHandler:   handler,
	}

	// Start the specified number of workers
	for range p.numOfWorkers {
		go p.startWorker()
	}

	return &p
}

// Submit adds a new job to the pool's job channel and increments the WaitGroup counter.
// The job is then processed by an available worker.
func (p *pool[T]) Submit(job T) {
	p.wg.Add(1)
	p.jobCh <- job
}

// StopAndWait stops the worker pool, waits for all jobs to be processed.
// It closes the job channel and waits for the workers to finish processing all submitted jobs.
func (p *pool[T]) StopAndWait() {
	close(p.jobCh)
	p.wg.Wait()
}

// startWorker defines every worker in pool.
// It continuously listens for jobs on the job channel and processes them using the job handler.
func (p *pool[T]) startWorker() {
	for {
		j := <-p.jobCh
		err := p.jobHandler(j)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		p.wg.Done()
	}
}
