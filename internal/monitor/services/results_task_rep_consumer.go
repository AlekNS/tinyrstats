package services

import (
	"context"

	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/runner"
)

type resultsTaskRepositoryConsumer struct {
	taskRepository monitor.TaskRepository
}

// Accept receives data and save it to the repository.
func (tc *resultsTaskRepositoryConsumer) Accept(ctx context.Context, results ...interface{}) error {
	// construct task from health task and status
	task := &monitor.Task{
		HealthTask: *results[0].(*monitor.HealthTask),
		Status:     results[1].(*monitor.HealthTaskStatus),
	}

	return tc.taskRepository.Save(ctx, task)
}

// NewResultsTaskRepositoryConsumer creates special task consumer to store it.
func NewResultsTaskRepositoryConsumer(taskRepository monitor.TaskRepository) runner.Consumer {
	return &resultsTaskRepositoryConsumer{
		taskRepository: taskRepository,
	}
}
