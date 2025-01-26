package workerpool

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type jobCh[T any] chan T
type jobHandler[T any] func(T) error

type pool[T any] struct {
	numOfWorkers int            // number of workers in pool
	jobCh        jobCh[T]       // chanel with jobs
	jobHandler   jobHandler[T]  // job handler
	wg           sync.WaitGroup // synchronizer of workers jobs
}

func New[T any](numOfWorkers, jobChSize int, handler jobHandler[T]) *pool[T] {
	jobCh := make(chan T, jobChSize)
	p := pool[T]{
		jobCh:        jobCh,
		numOfWorkers: numOfWorkers,
		jobHandler:   handler,
	}

	for range p.numOfWorkers {
		go p.startWorker()
	}

	return &p
}

// add a new job in worker pool
func (p *pool[T]) Submit(job T) {
	p.wg.Add(1)
	p.jobCh <- job
}

// stop and wait for the end of worker pool job handling
func (p *pool[T]) StopAndWait() {
	close(p.jobCh)
	p.wg.Wait()
}

// start new worker in pool
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
