package app

import (
	"context"
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	opentracing "github.com/opentracing/opentracing-go"

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

	level.Debug(logger).Log("msg", "create task", "url", req.Task.URL)

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

	level.Debug(logger).Log("msg", "schedule task")

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
	span, _ := opentracing.StartSpanFromContext(ctx, "QueryBy")
	defer span.Finish()

	logger := log.With(ta.logger, "method", "QueryBy")

	if len(query.ByHost) > 0 {
		level.Debug(logger).Log("msg", "query by resource")

		span.SetTag("type", "resource")

		task, err := ta.taskRepository.GetByID(ctx, monitor.TaskID(query.ByHost))
		if err != nil {
			span.LogEvent("not_found")
			return nil, err
		}

		span.SetTag("resource", string(task.ID))

		ta.events.TaskQueriedByResource().Emit(string(task.ID))

		span.LogEvent("response")

		return &monitor.QueryTaskResult{
			Task: *task,
		}, nil
	}

	isMaxQueried := query.ByResponseTime == monitor.QueryResponseMaxTime

	if isMaxQueried {
		level.Debug(logger).Log("msg", "query by max response time")

		span.SetTag("type", "max")
	} else {
		level.Debug(logger).Log("msg", "query by min response time")

		span.SetTag("type", "min")
	}

	task, err := ta.taskRepository.GetByResponseTimeMinOrMax(ctx, isMaxQueried)
	if err != nil {
		span.LogEvent("not_found")
		return nil, err
	}

	span.SetTag("resource", string(task.ID))

	if isMaxQueried {
		ta.events.TaskQueriedByMaxResponse().Emit(string(task.ID))
	} else {
		ta.events.TaskQueriedByMinResponse().Emit(string(task.ID))
	}

	span.LogEvent("response")

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
