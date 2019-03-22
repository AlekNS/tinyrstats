package app

import (
	"context"
	"net/url"

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

	// @TODO: For simplicity, but should be used hash.
	parsedURL, err := url.Parse(req.Task.URL)
	if err != nil {
		return nil, err
	}

	req.Task.ID = monitor.TaskID(parsedURL.Hostname())
	taskID := req.Task.ID

	_, err = ta.taskRepository.GetByID(ctx, monitor.TaskID(taskID))
	if err == nil {
		level.Debug(logger).Log("msg", "task already created")
		return &monitor.CreateTaskResult{
			ID: string(taskID),
		}, nil
	}

	task := &monitor.HealthTask{
		ID:      taskID,
		Timeout: req.Task.Timeout,
		Method:  req.Task.Method,
		URL:     req.Task.URL,
		Body:    req.Task.Body,
		Headers: req.Task.Headers,
	}

	if err := ta.scheduler.Schedule(ctx, taskID, &monitor.ScheduleHealthTask{
		Interval: 0,
		Task:     task,
	}); err != nil {
		return nil, err
	}

	return &monitor.CreateTaskResult{
		ID: string(taskID),
	}, nil
}

// GetStatus .
func (ta *taskAppImpl) QueryBy(ctx context.Context, query *monitor.QueryTask) (*monitor.QueryTaskResult, error) {
	if len(query.ByHost) > 0 {
		task, err := ta.taskRepository.GetByID(ctx, monitor.TaskID(query.ByHost))
		if err != nil {
			return nil, err
		}

		ta.events.TaskQueriedByURL().Emit(string(task.ID))

		return &monitor.QueryTaskResult{
			Task: *task,
		}, nil
	}

	task, err := ta.taskRepository.GetByResponseTimeMinOrMax(ctx, query.ByResponseTime == monitor.QueryResponseMaxTime)
	if err != nil {
		return nil, err
	}

	if query.ByResponseTime == monitor.QueryResponseMaxTime {
		ta.events.TaskQueriedByMaxResponse().Emit(string(task.ID))
	} else {
		ta.events.TaskQueriedByMinResponse().Emit(string(task.ID))
	}

	return &monitor.QueryTaskResult{
		Task: *task,
	}, nil
}

// newTaskApp .
func newTaskApp(
	settings *config.Settings,
	logger log.Logger,
	events monitor.Events,
	scheduler monitor.ScheduleTaskService,
	taskRepository monitor.TaskRepository) *taskAppImpl {

	return &taskAppImpl{
		logger:         log.With(logger, "service", "TaskApp"),
		settings:       settings,
		events:         events,
		taskRepository: taskRepository,
		scheduler:      scheduler,
	}
}
