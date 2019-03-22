package services

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/runner"
	"github.com/alekns/tinyrstats/pkg/helpers/str"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type (
	intervalScheduleItem struct {
		id string
		// interval in seconds
		interval int
		// scheduleAt is schedule time label of the task in future
		scheduleAt int
		// task is abstract object for passing it to the consumer
		task interface{}
	}

	intervalScheduleItems []*intervalScheduleItem

	// IntervalScheduler is used to schedule task periodically.
	// Each task can have an individual interval.
	IntervalScheduler struct {
		// consumer receives a scheduled task
		consumer runner.Consumer

		settings *config.SchedulerSettings

		logger log.Logger

		mtxs []sync.Mutex
		// items are splitted by smaller heaps for preventing long locks
		items []*intervalScheduleItems

		stopScheduleCh chan struct{}
		waitWorkers    sync.WaitGroup
	}
)

func getNextScheduleTime(interval int) int {
	return int(time.Now().Unix()) + interval
}

func (si intervalScheduleItems) Len() int { return len(si) }

func (si intervalScheduleItems) Less(i, j int) bool {
	return si[i].scheduleAt > si[j].scheduleAt
}

func (si intervalScheduleItems) Swap(i, j int) {
	si[i], si[j] = si[j], si[i]
}

func (si *intervalScheduleItems) Push(x interface{}) {
	item := x.(*intervalScheduleItem)
	*si = append(*si, item)
}

func (si *intervalScheduleItems) Pop() interface{} {
	old := *si
	n := len(old)
	item := old[n-1]
	*si = old[0 : n-1]
	return item
}

func (si *intervalScheduleItems) Top() interface{} {
	return (*si)[si.Len()-1]
}

func (si *intervalScheduleItems) FindIndex(id string) int {
	for i, item := range *si {
		if item.id == id {
			return i
		}
	}
	return -1
}

// Start starts workers.
func (st *IntervalScheduler) Start(ctx context.Context) error {
	st.waitWorkers.Add(st.settings.MaxConcurrency)

	for i := 0; i < st.settings.MaxConcurrency; i++ {
		go func(workerIndex int) {
			logger := log.With(st.logger, "worker", workerIndex)
			timer := time.NewTimer(0)
		stopSchedule:
			for {
				select {
				case <-st.stopScheduleCh:
					break stopSchedule
				case <-ctx.Done():

				case <-timer.C:
					items := make(intervalScheduleItems, 0)

					st.mtxs[workerIndex].Lock()

					nowAt := int(time.Now().Unix())

					// Get all reaching items
					workerItems := st.items[workerIndex]
					for workerItems.Len() > 0 && workerItems.Top().(*intervalScheduleItem).scheduleAt <= nowAt {
						items = append(items, workerItems.Pop().(*intervalScheduleItem))
					}

					// Reschedule these items again
					for _, item := range items {
						item.scheduleAt = getNextScheduleTime(item.interval)
						heap.Push(workerItems, item)
					}
					st.mtxs[workerIndex].Unlock()

					// Send all reached items
					for _, item := range items {
						if err := st.consumer.Accept(ctx, item.task); err != nil {
							// @TODO: use a backoff method instead of simple log!?
							level.Error(logger).Log("err", err.Error())
						}
					}

					timer.Reset(time.Second)
				}
			}

			timer.Stop()
			st.waitWorkers.Done()
		}(i)
	}

	return nil
}

// Stop and wait the shutdown of workers.
func (st *IntervalScheduler) Stop(ctx context.Context) error {
	close(st.stopScheduleCh)
	st.waitWorkers.Wait()
	return nil
}

// Schedule enqueue of task.
func (st *IntervalScheduler) Schedule(ctx context.Context, taskID monitor.TaskID, args *monitor.ScheduleHealthTask) error {
	logger := log.With(st.logger, "method", "Schedule", "taskId", taskID)

	workerIndex := str.BasicStrHash(string(taskID)) % st.settings.MaxConcurrency

	interval := args.Interval
	if interval < 1 {
		interval = st.settings.DefaultInterval
	}

	st.mtxs[workerIndex].Lock()

	// remove previous task with the same id to be idempotent.
	st.remove(args.Task.ID, workerIndex)

	heap.Push(st.items[workerIndex], &intervalScheduleItem{
		id:         string(args.Task.ID),
		scheduleAt: getNextScheduleTime(0), // @TODO: Is need to run task immediately!?
		interval:   interval,
		task:       args.Task,
	})

	st.mtxs[workerIndex].Unlock()

	level.Debug(logger).Log("msg", "place task")

	return nil
}

// remove is method for stopping a schedule of the task
func (st *IntervalScheduler) remove(taskID monitor.TaskID, workerIndex int) error {
	for inx, item := range *st.items[workerIndex] {
		if monitor.TaskID(item.id) == taskID {
			heap.Remove(st.items[workerIndex], inx)
			return nil
		}
	}

	return monitor.ErrTaskNotFound
}

// Cancel stops a task.
func (st *IntervalScheduler) Cancel(ctx context.Context, taskID monitor.TaskID) error {
	logger := log.With(st.logger, "method", "Cancel", "taskId", taskID)

	workerIndex := str.BasicStrHash(string(taskID)) % st.settings.MaxConcurrency

	st.mtxs[workerIndex].Lock()
	defer st.mtxs[workerIndex].Unlock()

	level.Debug(logger).Log("msg", "cancel task")

	return st.remove(taskID, workerIndex)
}

// CancelAll stops all tasks.
func (st *IntervalScheduler) CancelAll(ctx context.Context) error {
	logger := log.With(st.logger, "method", "CancelAll")

	for i := 0; i < st.settings.MaxConcurrency; i++ {
		st.mtxs[i].Lock()
		defer st.mtxs[i].Unlock()
	}

	st.init()

	level.Debug(logger).Log("msg", "complete")

	return nil
}

// init clean initialization.
func (st *IntervalScheduler) init() {
	st.mtxs = make([]sync.Mutex, st.settings.MaxConcurrency)
	st.items = make([]*intervalScheduleItems, st.settings.MaxConcurrency)
	for i := 0; i < st.settings.MaxConcurrency; i++ {
		st.items[i] = new(intervalScheduleItems)
	}
}

// NewIntervalScheduler creates interval scheduler.
func NewIntervalScheduler(logger log.Logger,
	settings *config.SchedulerSettings,
	consumer runner.Consumer) *IntervalScheduler {

	scheduler := &IntervalScheduler{
		logger:         logger,
		consumer:       consumer,
		settings:       settings,
		stopScheduleCh: make(chan struct{}),
	}

	scheduler.init()

	return scheduler
}
