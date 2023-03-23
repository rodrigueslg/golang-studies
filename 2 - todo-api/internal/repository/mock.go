package repository

import (
	"gochallenges/internal/model"
	"reflect"

	"github.com/golang/mock/gomock"
)

type TaskMock struct {
	ctrl     *gomock.Controller
	recorder *TaskMockMockRecorder
}

type TaskMockMockRecorder struct {
	mock *TaskMock
}

func NewTaskMock(ctrl *gomock.Controller) *TaskMock {
	mock := &TaskMock{ctrl: ctrl}
	mock.recorder = &TaskMockMockRecorder{mock}
	return mock
}

func (m *TaskMock) EXPECT() *TaskMockMockRecorder {
	return m.recorder
}

func (m *TaskMock) Create(task model.Task) (model.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", task)
	ret0, _ := ret[0].(model.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TaskMockMockRecorder) Create(task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*TaskMock)(nil).Create), task)
}

func (m *TaskMock) FindAll() ([]model.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll")
	ret0, _ := ret[0].([]model.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TaskMockMockRecorder) FindAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*TaskMock)(nil).FindAll))
}

func (m *TaskMock) FindByID(id int) (model.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(model.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TaskMockMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*TaskMock)(nil).FindByID), id)
}

func (m *TaskMock) FindByStatus(completed bool) ([]model.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByStatus", completed)
	ret0, _ := ret[0].([]model.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TaskMockMockRecorder) FindByStatus(completed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByStatus", reflect.TypeOf((*TaskMock)(nil).FindByStatus), completed)
}

func (m *TaskMock) Update(task model.Task) (model.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", task)
	ret0, _ := ret[0].(model.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *TaskMockMockRecorder) Update(task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*TaskMock)(nil).Update), task)
}

func (m *TaskMock) Delete(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *TaskMockMockRecorder) Delete(id int) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*TaskMock)(nil).Delete), id)
}

func (m *TaskMock) Close() {
}
