package monitor

//go:generate mockgen -source=interfaces.go -package=monitor -destination=interfaces_mocks.go

import (
	"context"

	"github.com/alekns/tinyrstats/pkg/helpers/subscribs"
)

type (
	// TaskRepository is a storage for tasks.
	TaskRepository interface {
		// GetByResponseTimeMinOrMax returns task with most minimal or maximal response time.
		GetByResponseTimeMinOrMax(context.Context, bool) (*Task, error)
		// GetByID returns task by identifier.
		GetByID(context.Context, TaskID) (*Task, error)

		// Save puts into storage.
		Save(context.Context, *Task) error

		// Delete by identifier.
		Delete(context.Context, TaskID) error
		// DeleteAll cleans repository.
		DeleteAll(context.Context)
	}

	// HealthService is a health checker.
	HealthService interface {
		// CheckStatus returns status of the resource.
		CheckStatus(context.Context, *HealthTask) (*HealthTaskStatus, error)
	}

	// ScheduleTaskService is need for run tasks periodically.
	ScheduleTaskService interface {
		// Schedule setups task in a schedule.
		Schedule(context.Context, TaskID, *ScheduleHealthTask) error
		// Cancel a task by identifier.
		Cancel(context.Context, TaskID) error
		// CancelAll cleans a schedule.
		CancelAll(context.Context) error
	}

	// Events all processable events.
	Events interface {
		// TaskQueriedByResource is fired when task was quired by resource.
		TaskQueriedByResource() subscribs.EventHandler
		// TaskQueriedByMinResponse is fired when task was quired.
		TaskQueriedByMinResponse() subscribs.EventHandler
		// TaskQueriedByMaxResponse is fired when task was quired.
		TaskQueriedByMaxResponse() subscribs.EventHandler
	}

	// StatsHostsInfo holds call count by resources.
	StatsHostsInfo = map[string]int

	// StatsService gathering statistics of the requests
	StatsService interface {
		// GetAllHosts returns statistics by all resources.
		GetAllHosts() StatsHostsInfo
		// GetMinMax returns min and max requests count.
		GetMinMax() (int32, int32)

		// AddHost changes counter for resource.
		AddHost(string, int)
		// AddMinMax changes counter for min or max.
		AddMinMax(bool, int32)

		// DeleteHost removes resource counter.
		DeleteHost(string)
	}

	// Applications (singletons)

	// TaskApp provides complete logic for working with a task.
	TaskApp interface {
		// CreateAndRun registers and schedules a task.
		CreateAndRun(context.Context, *CreateTaskCommand) (*CreateTaskResult, error)
		// QueryBy gets task by parameters.
		QueryBy(context.Context, *QueryTask) (*QueryTaskResult, error)
	}

	// StatsApp provides statistics of count TaskApp.QueryBy.
	StatsApp interface {
		// QueryBy gets stats by parameters.
		QueryBy(context.Context, *QueryCallStatistic) (*QueryCallStatisticResult, error)
	}
)
