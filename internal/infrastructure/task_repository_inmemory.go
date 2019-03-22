package infrastructure

import (
	"context"
	"sync"
	"sync/atomic"

	uuid "github.com/satori/go.uuid"

	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/str"
)

// @TODO: Use from config
const preDeclaredDataSize = 64

type (
	taskByIDBucket = map[monitor.TaskID]*monitor.Task

	taskRepositoryInMemoryImpl struct {
		bucketsCount int

		mtxs    []sync.RWMutex
		buckets []taskByIDBucket

		minResponseTask atomic.Value
		maxResponseTask atomic.Value
	}
)

var (
	nilResponseTask = &monitor.Task{}
)

// cloneTask not safe (skip headers)
func cloneTask(src *monitor.Task) *monitor.Task {
	dst := &monitor.Task{}
	*dst = *src
	if src.Status != nil {
		dst.Status = &monitor.HealthTaskStatus{}
		// @TODO: +Copy Header
		*dst.Status = *src.Status
	}
	return dst
}

// GetByResponseTimeMinOrMax .
func (tr *taskRepositoryInMemoryImpl) GetByResponseTimeMinOrMax(ctx context.Context, isNeedMax bool) (*monitor.Task, error) {
	var task *monitor.Task

	if isNeedMax {
		task = tr.maxResponseTask.Load().(*monitor.Task)
	} else {
		task = tr.minResponseTask.Load().(*monitor.Task)
	}

	if task == nilResponseTask {
		return nil, monitor.ErrTaskNotFound
	}

	return cloneTask(task), nil
}

// GetByID .
func (tr *taskRepositoryInMemoryImpl) GetByID(ctx context.Context, taskID monitor.TaskID) (*monitor.Task, error) {
	bucketIndex := str.BasicStrHash(string(taskID)) % tr.bucketsCount

	tr.mtxs[bucketIndex].RLock()
	defer tr.mtxs[bucketIndex].RUnlock()

	if elem, ok := tr.buckets[bucketIndex][taskID]; ok {
		return cloneTask(elem), nil
	}

	return nil, monitor.ErrTaskNotFound
}

// Save .
func (tr *taskRepositoryInMemoryImpl) Save(ctx context.Context, task *monitor.Task) error {
	bucketIndex := str.BasicStrHash(string(task.ID)) % tr.bucketsCount

	if len(task.ID) == 0 {
		task.ID = monitor.TaskID(uuid.NewV4().String())
	}

	task = cloneTask(task)

	tr.mtxs[bucketIndex].Lock()
	tr.buckets[bucketIndex][task.ID] = task
	tr.mtxs[bucketIndex].Unlock()

	tr.updateMinMaxResponseTasks(task, true)

	return nil
}

// updateMinMaxResponseTasks update min and max response tasks after update or delete.
func (tr *taskRepositoryInMemoryImpl) updateMinMaxResponseTasks(task *monitor.Task, isUpdate bool) {
	if task.Status == nil {
		return
	}

	if isUpdate {
		minTask := tr.minResponseTask.Load().(*monitor.Task)
		maxTask := tr.maxResponseTask.Load().(*monitor.Task)
		if minTask == nilResponseTask {
			tr.minResponseTask.Store(task)
		}
		if maxTask == nilResponseTask {
			tr.maxResponseTask.Store(task)
		}
		if minTask.Status != nil && minTask.Status.ResponseTime > task.Status.ResponseTime {
			tr.minResponseTask.Store(task)
		}
		if maxTask.Status != nil && maxTask.Status.ResponseTime < task.Status.ResponseTime {
			tr.maxResponseTask.Store(task)
		}
		return
	}

	// Ugly part. After deletion need to update minResponseTask and maxResponseTask
	// by iteration over items of buckets or better to use another data structures
	// like binary trees instead of a simple hash map (may be).
	minTask := tr.minResponseTask.Load().(*monitor.Task)
	maxTask := tr.maxResponseTask.Load().(*monitor.Task)
	statTaskIndex := 0

	if minTask == task {
		statTaskIndex = 1
		minTask = nilResponseTask
	} else if maxTask == task {
		statTaskIndex = 2
		maxTask = nilResponseTask
	}

	if statTaskIndex > 0 {
		for _, bucket := range tr.buckets {
			for _, task := range bucket {
				if task.Status == nil {
					continue
				}
				if statTaskIndex == 1 && (minTask == nilResponseTask || task.Status.ResponseTime < minTask.Status.ResponseTime) {
					minTask = task
				} else if statTaskIndex == 2 && (maxTask == nilResponseTask || task.Status.ResponseTime > maxTask.Status.ResponseTime) {
					maxTask = task
				}
			}
		}
		if statTaskIndex == 1 {
			tr.minResponseTask.Store(minTask)
		} else if statTaskIndex == 2 {
			tr.maxResponseTask.Store(maxTask)
		}
	}

}

// Delete .
func (tr *taskRepositoryInMemoryImpl) Delete(ctx context.Context, taskID monitor.TaskID) error {
	bucketIndex := str.BasicStrHash(string(taskID)) % tr.bucketsCount

	tr.mtxs[bucketIndex].Lock()

	if task, ok := tr.buckets[bucketIndex][taskID]; ok {
		delete(tr.buckets[bucketIndex], taskID)
		tr.mtxs[bucketIndex].Unlock()

		tr.updateMinMaxResponseTasks(task, false)

		return nil
	}
	tr.mtxs[bucketIndex].Unlock()

	return monitor.ErrTaskNotFound
}

// DeleteAll .
func (tr *taskRepositoryInMemoryImpl) DeleteAll(ctx context.Context) {
	for i := 0; i < tr.bucketsCount; i++ {
		tr.mtxs[i].Lock()
		defer tr.mtxs[i].Unlock()
	}

	initTaskRepository(tr)
}

func initTaskRepository(tr *taskRepositoryInMemoryImpl) *taskRepositoryInMemoryImpl {
	tr.mtxs = make([]sync.RWMutex, tr.bucketsCount)
	tr.buckets = make([]taskByIDBucket, tr.bucketsCount)
	for i := 0; i < tr.bucketsCount; i++ {
		tr.buckets[i] = make(taskByIDBucket, preDeclaredDataSize)
	}

	tr.minResponseTask.Store(nilResponseTask)
	tr.maxResponseTask.Store(nilResponseTask)

	return tr
}

// NewTaskRepositoryInMemory .
func NewTaskRepositoryInMemory(bucketsCount int) monitor.TaskRepository {
	return initTaskRepository(&taskRepositoryInMemoryImpl{
		bucketsCount: bucketsCount,
	})
}
