package controller

import "github.com/golang/mock/gomock"

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
