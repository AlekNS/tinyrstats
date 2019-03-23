package tasks

import (
	"context"

	"github.com/alekns/tinyrstats/internal/monitor"
)

// ScheduleTasksSlice gets slice of tasks and launch it.
func ScheduleTasksSlice(ctx context.Context,
	taskApp monitor.TaskApp, items []*monitor.ScheduleHealthTask) error {

	var err error
	for _, item := range items {
		_, err = taskApp.CreateAndRun(ctx, &monitor.CreateTaskCommand{ScheduleHealthTask: *item})
		if err != nil {
			return err
		}
	}

	return nil
}
