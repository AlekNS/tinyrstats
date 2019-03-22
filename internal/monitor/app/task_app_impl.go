package app

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
)

type taskAppImpl struct {
	settings *config.Settings
	logger   log.Logger

	events         monitor.Events
	taskRepository monitor.TaskRepository
	scheduler      monitor.ScheduleTaskService
}

// Create .
func (ta *taskAppImpl) CreateAndRun(ctx context.Context, req *monitor.CreateTaskCommand) (*monitor.CreateTaskResult, error) {
	logger := log.With(ta.logger, "method", "Create")

	level.Debug(logger).Log("msg", "CreateAndRun")

	// taskID := req.URL

	// task, err := ta.taskRepository.GetByID(taskID)
	// if err == nil {
	// 	// Check if task is already in queue
	// 	if task.Status.Status == 0 && task.Status.IsPending {
	// 		logger.Debug("task already in pending state", task.Request)
	// 		return &monitor.CreateTaskResponse{
	// 			ID: taskID,
	// 		}, nil
	// 	}
	// }

	// task = &monitor.Task{
	// 	ID: taskID,
	// 	HealthTask: monitor.HealthTask{
	// 		Timeout: req.Timeout,
	// 		Method:  req.Method,
	// 		URL:     req.URL,
	// 		Body:    req.Body,
	// 		Headers: req.Headers,
	// 	},
	// }

	// _, err = ta.requestResource.EnqueueTask(ctx, task)
	// if err != nil {
	// 	return nil, err
	// }

	return &monitor.CreateTaskResult{}, nil
}

// GetStatus .
func (ta *taskAppImpl) QueryBy(ctx context.Context, query *monitor.QueryTask) (*monitor.QueryTaskResult, error) {
	if len(query.ByHost) > 0 {
		ta.events.TaskQueriedByURL().Emit(query.ByHost)
	}

	// ta.events.TaskQueriedByMinResponse(task.URL)
	// ta.events.TaskQueriedByMaxResponse(task.URL)

	// task, err := ta.taskRepository.GetByID(req.ID)
	// if err != nil {
	// 	return nil, err
	// }

	return &monitor.QueryTaskResult{}, nil
}

// newTaskApp .
func newTaskApp(
	settings *config.Settings,
	logger log.Logger,
	taskRepository monitor.TaskRepository) *taskAppImpl {

	return &taskAppImpl{
		logger:         log.With(logger, "service", "TaskApp"),
		settings:       settings,
		taskRepository: taskRepository,
	}
}
