package services

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alekns/tinyrstats/internal/monitor"
)

func getPreFilledTaskRepository(cnt int) monitor.TaskRepository {
	rep := NewTaskRepositoryInMemory(4)

	var t *monitor.Task

	for i := 0; i < cnt; i++ {
		t = &monitor.Task{
			HealthTask: monitor.HealthTask{ID: monitor.TaskID(strconv.Itoa(i))},
			Status: &monitor.HealthTaskStatus{
				ResponseTime: int64(i + 1),
			},
		}
		if i&1 == 1 {
			t.Status = nil
		}
		rep.Save(context.TODO(), t)
	}

	return rep
}

func TestInMemTaskRepositoryShouldSuccessGetById(t *testing.T) {
	rep := getPreFilledTaskRepository(8)

	result, err := rep.GetByID(context.TODO(), "2")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, monitor.TaskID("2"), result.ID)
}

func TestInMemTaskRepositoryShouldNotFoundGetById(t *testing.T) {
	rep := getPreFilledTaskRepository(8)

	result, err := rep.GetByID(context.TODO(), monitor.TaskID("invalidId"))

	require.Error(t, err, monitor.ErrTaskNotFound)
	require.Nil(t, result)
}

func TestInMemTaskRepositoryShouldUpdateTask(t *testing.T) {
	rep := getPreFilledTaskRepository(8)
	task, err := rep.GetByID(context.TODO(), "1")

	require.NoError(t, err, "expect no error when getting by id 1")
	require.NotNil(t, task)

	task.Method = "GET"
	err = rep.Save(context.TODO(), task)

	require.NoError(t, err, "expect no error after save")

	// To check that it's not pointer
	task.Method = "POST"

	task, err = rep.GetByID(context.TODO(), "1")

	require.NoError(t, err, "expect no error after getting by id 1 after save")
	require.NotNil(t, task, "expect not nil task after getting by id 1 after save")

	require.NotNil(t, task)
	require.Equal(t, "GET", task.Method)
}

func TestInMemTaskRepositoryShouldDeleteTask(t *testing.T) {
	rep := getPreFilledTaskRepository(8)

	err := rep.Delete(context.TODO(), "1")

	require.NoError(t, err, "expect no error after deletion by id 1")

	result, err := rep.GetByID(context.TODO(), "1")

	require.Error(t, err, monitor.ErrTaskNotFound)
	require.Nil(t, result)
}

func TestInMemTaskRepositoryShouldDeleteAllTasks(t *testing.T) {
	rep := getPreFilledTaskRepository(8)

	rep.DeleteAll(context.TODO())

	result, err := rep.GetByID(context.TODO(), "5")

	require.Error(t, err, monitor.ErrTaskNotFound)
	require.Nil(t, result)
}

func TestInMemTaskRepositoryMinMaxOnEmpty(t *testing.T) {
	rep := NewTaskRepositoryInMemory(8)

	minTask, err := rep.GetByResponseTimeMinOrMax(context.TODO(), false)
	require.Error(t, err, monitor.ErrTaskNotFound)
	require.Nil(t, minTask)

	maxTask, err := rep.GetByResponseTimeMinOrMax(context.TODO(), true)
	require.Error(t, err, monitor.ErrTaskNotFound)
	require.Nil(t, maxTask)
}

func TestInMemTaskRepositoryMinMaxSuccess(t *testing.T) {
	rep := getPreFilledTaskRepository(8)

	minTask, err := rep.GetByResponseTimeMinOrMax(context.TODO(), false)
	require.NoError(t, err)
	require.NotNil(t, minTask)
	require.NotNil(t, minTask.Status)
	require.Equal(t, int64(1), minTask.Status.ResponseTime)

	maxTask, err := rep.GetByResponseTimeMinOrMax(context.TODO(), true)
	require.NoError(t, err)
	require.NotNil(t, maxTask)
	require.NotNil(t, maxTask.Status)
	require.Equal(t, int64(7), maxTask.Status.ResponseTime)
}

func TestInMemTaskRepositoryMinMaxAfterDeleteSuccess(t *testing.T) {
	rep := getPreFilledTaskRepository(8)

	rep.Delete(context.TODO(), "0")
	rep.Delete(context.TODO(), "7")

	minTask, err := rep.GetByResponseTimeMinOrMax(context.TODO(), false)
	require.NoError(t, err)
	require.NotNil(t, minTask)
	require.NotNil(t, minTask.Status)
	require.Equal(t, int64(3), minTask.Status.ResponseTime)

	maxTask, err := rep.GetByResponseTimeMinOrMax(context.TODO(), true)
	require.NoError(t, err)
	require.NotNil(t, maxTask)
	require.NotNil(t, maxTask.Status)
	require.Equal(t, int64(7), maxTask.Status.ResponseTime)
}
