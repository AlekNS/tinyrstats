package runner

import (
	"context"
	"sync"
	"sync/atomic"
)

// ConcurrentProcessor implements processing of tasks in concurrency mode.
type ConcurrentProcessor struct {
	mtx           sync.Mutex
	waitConsumers sync.WaitGroup

	enqueueLimiter int32
	enqueueCount   int32
	consumersCount int

	consumerQueue chan interface{}
	stopProcessor chan struct{}

	isStopped bool
}

// Stop processing of the tasks.
func (p *ConcurrentProcessor) Stop() error {
	p.mtx.Lock()
	if !p.isStopped {
		close(p.stopProcessor)
		p.isStopped = true
	}
	p.mtx.Unlock()

	return nil
}

// Enqueue place task to the queue.
func (p *ConcurrentProcessor) Enqueue(ctx context.Context, task interface{}) error {
	if p.enqueueLimiter > 0 {
		cnt := atomic.LoadInt32(&p.enqueueCount)
		if cnt > p.enqueueLimiter {
			return ErrBackpressure
		}

		atomic.AddInt32(&p.enqueueCount, 1)
	}

	select {
	case <-ctx.Done():
	case p.consumerQueue <- task:
	}

	return nil
}

// Start set of workers.
func (p *ConcurrentProcessor) Start(ctx context.Context, consumer Consumer, errors ErrorHandler) error {
	p.waitConsumers.Add(p.consumersCount)

	for worker := 0; worker < p.consumersCount; worker++ {
		go func() {
			var err error
		stopWorker:
			for {
				select {
				case <-p.stopProcessor:
					break stopWorker
				case <-ctx.Done():
					break stopWorker

				case task, ok := <-p.consumerQueue:
					if !ok {
						break stopWorker
					}

					err = consumer.Accept(ctx, task)
					if err != nil {
						err = errors.OnError(err)
						if err != nil {
							p.Stop()
						}
					}

					if p.enqueueLimiter > 0 {
						atomic.AddInt32(&p.enqueueCount, -1)
					}
				}
			}

			p.waitConsumers.Done()
		}()
	}

	return nil
}

// Wait until all workers will be stopped.
func (p *ConcurrentProcessor) Wait() {
	p.waitConsumers.Wait()
}

// NewConcurrentProcessor create concurrent processor.
func NewConcurrentProcessor(consumerQueueSize, consumersCount, enqueueLimiter int) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		enqueueLimiter: int32(enqueueLimiter),
		consumersCount: consumersCount,
		consumerQueue:  make(chan interface{}, consumerQueueSize),
		stopProcessor:  make(chan struct{}),
	}
}
