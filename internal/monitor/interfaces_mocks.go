// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package worker is a generated GoMock package.
package worker

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTaskRepository is a mock of TaskRepository interface
type MockTaskRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTaskRepositoryMockRecorder
}

// MockTaskRepositoryMockRecorder is the mock recorder for MockTaskRepository
type MockTaskRepositoryMockRecorder struct {
	mock *MockTaskRepository
}

// NewMockTaskRepository creates a new mock instance
func NewMockTaskRepository(ctrl *gomock.Controller) *MockTaskRepository {
	mock := &MockTaskRepository{ctrl: ctrl}
	mock.recorder = &MockTaskRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTaskRepository) EXPECT() *MockTaskRepositoryMockRecorder {
	return m.recorder
}

// GetByResponseTimeMinOrMax mocks base method
func (m *MockTaskRepository) GetByResponseTimeMinOrMax(ctx context.Context, bool isNeedMax) (*Task, error) {
	ret := m.ctrl.Call(m, "GetByResponseTimeMinOrMax", ctx, bool)
	ret0, _ := ret[0].(*Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByResponseTimeMinOrMax indicates an expected call of GetByResponseTimeMinOrMax
func (mr *MockTaskRepositoryMockRecorder) GetByResponseTimeMinOrMax(ctx, bool interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByResponseTimeMinOrMax", reflect.TypeOf((*MockTaskRepository)(nil).GetByResponseTimeMinOrMax), ctx, bool)
}

// GetByID mocks base method
func (m *MockTaskRepository) GetByID(ctx context.Context, taskID TaskID) (*Task, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, taskID)
	ret0, _ := ret[0].(*Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockTaskRepositoryMockRecorder) GetByID(ctx, taskID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTaskRepository)(nil).GetByID), ctx, taskID)
}

// Save mocks base method
func (m *MockTaskRepository) Save(ctx context.Context, task *Task) error {
	ret := m.ctrl.Call(m, "Save", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockTaskRepositoryMockRecorder) Save(ctx, task interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTaskRepository)(nil).Save), ctx, task)
}

// Delete mocks base method
func (m *MockTaskRepository) Delete(ctx context.Context, taskID TaskID) error {
	ret := m.ctrl.Call(m, "Delete", ctx, taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTaskRepositoryMockRecorder) Delete(ctx, taskID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTaskRepository)(nil).Delete), ctx, taskID)
}

// DeleteAll mocks base method
func (m *MockTaskRepository) DeleteAll(ctx context.Context) {
	m.ctrl.Call(m, "DeleteAll", ctx)
}

// DeleteAll indicates an expected call of DeleteAll
func (mr *MockTaskRepositoryMockRecorder) DeleteAll(ctx interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockTaskRepository)(nil).DeleteAll), ctx)
}

// MockHealthService is a mock of HealthService interface
type MockHealthService struct {
	ctrl     *gomock.Controller
	recorder *MockHealthServiceMockRecorder
}

// MockHealthServiceMockRecorder is the mock recorder for MockHealthService
type MockHealthServiceMockRecorder struct {
	mock *MockHealthService
}

// NewMockHealthService creates a new mock instance
func NewMockHealthService(ctrl *gomock.Controller) *MockHealthService {
	mock := &MockHealthService{ctrl: ctrl}
	mock.recorder = &MockHealthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHealthService) EXPECT() *MockHealthServiceMockRecorder {
	return m.recorder
}

// CheckStatus mocks base method
func (m *MockHealthService) CheckStatus(arg0 context.Context, arg1 *HealthTask) (*HealthTaskStatus, error) {
	ret := m.ctrl.Call(m, "CheckStatus", arg0, arg1)
	ret0, _ := ret[0].(*HealthTaskStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckStatus indicates an expected call of CheckStatus
func (mr *MockHealthServiceMockRecorder) CheckStatus(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckStatus", reflect.TypeOf((*MockHealthService)(nil).CheckStatus), arg0, arg1)
}

// MockScheduleTaskService is a mock of ScheduleTaskService interface
type MockScheduleTaskService struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleTaskServiceMockRecorder
}

// MockScheduleTaskServiceMockRecorder is the mock recorder for MockScheduleTaskService
type MockScheduleTaskServiceMockRecorder struct {
	mock *MockScheduleTaskService
}

// NewMockScheduleTaskService creates a new mock instance
func NewMockScheduleTaskService(ctrl *gomock.Controller) *MockScheduleTaskService {
	mock := &MockScheduleTaskService{ctrl: ctrl}
	mock.recorder = &MockScheduleTaskServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScheduleTaskService) EXPECT() *MockScheduleTaskServiceMockRecorder {
	return m.recorder
}

// Schedule mocks base method
func (m *MockScheduleTaskService) Schedule(arg0 context.Context, arg1 *ScheduleHealthTask) error {
	ret := m.ctrl.Call(m, "Schedule", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Schedule indicates an expected call of Schedule
func (mr *MockScheduleTaskServiceMockRecorder) Schedule(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Schedule", reflect.TypeOf((*MockScheduleTaskService)(nil).Schedule), arg0, arg1)
}

// Cancel mocks base method
func (m *MockScheduleTaskService) Cancel(arg0 context.Context, arg1 TaskID) error {
	ret := m.ctrl.Call(m, "Cancel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cancel indicates an expected call of Cancel
func (mr *MockScheduleTaskServiceMockRecorder) Cancel(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockScheduleTaskService)(nil).Cancel), arg0, arg1)
}

// CancelAll mocks base method
func (m *MockScheduleTaskService) CancelAll(arg0 context.Context) error {
	ret := m.ctrl.Call(m, "CancelAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelAll indicates an expected call of CancelAll
func (mr *MockScheduleTaskServiceMockRecorder) CancelAll(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelAll", reflect.TypeOf((*MockScheduleTaskService)(nil).CancelAll), arg0)
}