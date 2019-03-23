package runner

import (
	"context"
	"sync"
	"sync/atomic"
)

// ConcurrentProcessor implements processing of tasks in concurrency mode.
type ConcurrentProcessor struct {
	mtx               sync.Mutex
	waitConsumersStop sync.WaitGroup

	// enqueueLimiter using for backpressure behavior
	enqueueLimiter int32
	// enqueueLimiter specified count of current tasks
	enqueueCount int32

	// consumersCount is goroutine count
	consumersCount int
	// consumerQueueSize is channel buffer size
	consumerQueueSize int

	consumerQueueCh chan interface{}
	stopProcessorCh chan struct{}
	isStopped       bool
}

// Stop processing of the tasks.
func (p *ConcurrentProcessor) Stop() error {
	p.mtx.Lock()
	if !p.isStopped {
		close(p.stopProcessorCh)
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
		return ctx.Err()
	case p.consumerQueueCh <- task:
	}

	return nil
}

// Start set of workers.
func (p *ConcurrentProcessor) Start(ctx context.Context, consumer Consumer, errors ErrorHandler) error {
	p.consumerQueueCh = make(chan interface{}, p.consumerQueueSize)
	p.stopProcessorCh = make(chan struct{})
	p.isStopped = false

	p.waitConsumersStop.Add(p.consumersCount)

	for worker := 0; worker < p.consumersCount; worker++ {
		go func() {
			var err error
		stopWorker:
			for {
				select {
				case <-p.stopProcessorCh:
					break stopWorker
				case <-ctx.Done():
					break stopWorker

				case task, ok := <-p.consumerQueueCh:
					if !ok {
						break stopWorker
					}

					err = consumer.Accept(ctx, task)
					if err != nil {
						if errors.OnError(err) != nil {
							// stop processor if any error was propagated
							p.Stop()
						}
					}

					if p.enqueueLimiter > 0 {
						atomic.AddInt32(&p.enqueueCount, -1)
					}
				}
			}

			p.waitConsumersStop.Done()
		}()
	}

	return nil
}

// Wait until all workers will be stopped.
func (p *ConcurrentProcessor) Wait() {
	p.waitConsumersStop.Wait()
}

// NewConcurrentProcessor create concurrent processor.
func NewConcurrentProcessor(consumerQueueSize, consumersCount, enqueueLimiter int) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		enqueueLimiter:    int32(enqueueLimiter),
		consumersCount:    consumersCount,
		consumerQueueSize: consumerQueueSize,
	}
}
