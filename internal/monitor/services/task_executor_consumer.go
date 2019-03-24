package services

import (
	"context"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/pkg/helpers/runner"
	"github.com/go-kit/kit/log"
)

type (
	// ConcurrentTaskExecutor processing a tasks in concurrent mode.
	ConcurrentTaskExecutor struct {
		logger    log.Logger
		processor *runner.ConcurrentProcessor

		taskResultConsumer runner.Consumer
		errorsHandler      runner.ErrorHandler
	}
)

// Accept task from producer and enqueue it.
func (ct *ConcurrentTaskExecutor) Accept(ctx context.Context, results ...interface{}) error {
	return ct.processor.Enqueue(ctx, results[0])
}

// Start workers.
func (ct *ConcurrentTaskExecutor) Start(ctx context.Context) error {
	return ct.processor.Start(ctx, ct.taskResultConsumer, ct.errorsHandler)
}

// Stop workers.
func (ct *ConcurrentTaskExecutor) Stop(ctx context.Context) error {
	if err := ct.processor.Stop(); err != nil {
		return err
	}

	ct.processor.Wait()
	return nil
}

// NewConcurretTaskExecutor creates concurrent task processor.
func NewConcurretTaskExecutor(settings *config.TasksSettings,
	logger log.Logger,
	taskResultConsumer runner.Consumer) *ConcurrentTaskExecutor {

	const svcName = "ConcurrentTaskExecutor"

	return &ConcurrentTaskExecutor{
		logger: log.With(logger, "service", svcName),
		processor: runner.NewConcurrentProcessor(
			settings.QueueSize, settings.MaxConcurrency, settings.MaxPending),

		taskResultConsumer: taskResultConsumer,
		errorsHandler:      NewLoggerRunnerException(logger, svcName),
	}
}
