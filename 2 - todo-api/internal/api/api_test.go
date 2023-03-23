package api_test

import (
	"encoding/json"
	"fmt"
	"gochallenges/internal/controller"
	"gochallenges/internal/model"
	"gochallenges/internal/repository"
	"gochallenges/pkg"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repository.NewTaskMock(ctrl)
	ctrlMock := controller.NewTask(repoMock)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

	r.Header.Set("Authorization", "wrongtoken")
	ctrlMock.ServeHTTP(w, r)

	assertStatusCode(t, w.Result().StatusCode, http.StatusUnauthorized)
}

func TestGetAll(t *testing.T) {
	tasks := []model.Task{{Name: "task mock", Completed: true}}
	cases := []struct {
		caseName           string
		expectedError      error
		expectedStatusCode int
		expectedBehavior   func(m *repository.TaskMock)
		expectedBody       []model.Task
	}{
		{
			caseName:           "successfully retrieve all tasks",
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindAll().DoAndReturn(func() ([]model.Task, error) {
					return tasks, nil
				}).Times(1)
			},
			expectedBody: tasks,
		},
		{
			caseName:           "error retrieving all tasks",
			expectedError:      model.ErrExecuteQuery,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindAll().DoAndReturn(func() ([]model.Task, error) {
					return nil, model.ErrExecuteQuery
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed scanning rows",
			expectedError:      model.ErrScanningRows,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindAll().DoAndReturn(func() ([]model.Task, error) {
					return nil, model.ErrScanningRows
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed connect to database",
			expectedError:      model.ErrConnectDatabase,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindAll().DoAndReturn(func() ([]model.Task, error) {
					return nil, model.ErrConnectDatabase
				}).Times(1)
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := repository.NewTaskMock(ctrl)
			ctrlMock := controller.NewTask(repoMock)

			testCase.expectedBehavior(repoMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
			r.Header.Set("Authorization", pkg.GetBearerToken())

			ctrlMock.ServeHTTP(w, r)

			if testCase.expectedError != nil {
				errorMessage, _ := ioutil.ReadAll(w.Body)
				assertError(t, string(errorMessage), testCase.expectedError.Error())
				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
			} else {
				var got []model.Task
				json.NewDecoder(w.Body).Decode(&got)

				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
				assertResponseBody(t, got, testCase.expectedBody)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	task := model.Task{ID: 1, Name: "task mock", Completed: true}
	cases := []struct {
		caseName           string
		expectedError      error
		expectedStatusCode int
		expectedBehavior   func(m *repository.TaskMock)
		expectedBody       model.Task
	}{
		{
			caseName:           "successfully retrieve a task",
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
			},
			expectedBody: task,
		},
		{
			caseName:           "error retrieving a task",
			expectedError:      model.ErrExecuteQuery,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, model.ErrExecuteQuery
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed scanning row",
			expectedError:      model.ErrScanningRows,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, model.ErrScanningRows
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed connect to database",
			expectedError:      model.ErrConnectDatabase,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, model.ErrConnectDatabase
				}).Times(1)
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := repository.NewTaskMock(ctrl)
			ctrlMock := controller.NewTask(repoMock)

			testCase.expectedBehavior(repoMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", task.ID), nil)
			r.Header.Set("Authorization", pkg.GetBearerToken())

			ctrlMock.ServeHTTP(w, r)

			if testCase.expectedError != nil {
				errorMessage, _ := ioutil.ReadAll(w.Body)
				assertError(t, string(errorMessage), testCase.expectedError.Error())
				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
			} else {
				var got model.Task
				json.NewDecoder(w.Body).Decode(&got)

				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
				assertResponseBody(t, got, testCase.expectedBody)
			}
		})
	}
}

func TestGetByStatus(t *testing.T) {
	completedTasks := []model.Task{{Name: "completed task mock", Completed: true}}
	uncomplemtedTasks := []model.Task{{Name: "uncompleted task mock", Completed: false}}
	cases := []struct {
		caseName           string
		expectedError      error
		expectedStatusCode int
		expectedBehavior   func(m *repository.TaskMock)
		expectedBody       []model.Task
		expectedStatus     bool
	}{
		{
			caseName:           "successfully retrieve completed",
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByStatus(true).DoAndReturn(func(completed bool) ([]model.Task, error) {
					return completedTasks, nil
				}).Times(1)
			},
			expectedBody:   completedTasks,
			expectedStatus: true,
		},
		{
			caseName:           "successfully retrieve uncompleted",
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByStatus(false).DoAndReturn(func(completed bool) ([]model.Task, error) {
					return uncomplemtedTasks, nil
				}).Times(1)
			},
			expectedBody:   uncomplemtedTasks,
			expectedStatus: false,
		},
		{
			caseName:           "error retrieving all tasks by status",
			expectedError:      model.ErrExecuteQuery,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByStatus(false).DoAndReturn(func(completed bool) ([]model.Task, error) {
					return nil, model.ErrExecuteQuery
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed scanning rows",
			expectedError:      model.ErrScanningRows,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByStatus(false).DoAndReturn(func(completed bool) ([]model.Task, error) {
					return nil, model.ErrScanningRows
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed connect to database",
			expectedError:      model.ErrConnectDatabase,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByStatus(false).DoAndReturn(func(completed bool) ([]model.Task, error) {
					return nil, model.ErrConnectDatabase
				}).Times(1)
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := repository.NewTaskMock(ctrl)
			ctrlMock := controller.NewTask(repoMock)

			testCase.expectedBehavior(repoMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/tasks?completed=%t", testCase.expectedStatus), nil)
			r.Header.Set("Authorization", pkg.GetBearerToken())

			ctrlMock.ServeHTTP(w, r)

			if testCase.expectedError != nil {
				errorMessage, _ := ioutil.ReadAll(w.Body)
				assertError(t, string(errorMessage), testCase.expectedError.Error())
				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
			} else {
				var got []model.Task
				json.NewDecoder(w.Body).Decode(&got)

				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
				assertResponseBody(t, got, testCase.expectedBody)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	task := model.Task{Name: "study golang unit testing", Completed: false}
	cases := []struct {
		caseName           string
		expectedError      error
		expectedStatusCode int
		expectedBehavior   func(m *repository.TaskMock)
		expectedBody       model.Task
		requestBody        string
	}{
		{
			caseName:           "successfully creating a task",
			expectedError:      nil,
			expectedStatusCode: http.StatusCreated,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().Create(task).DoAndReturn(func(newTask model.Task) (model.Task, error) {
					return task, nil
				}).Times(1)
			},
			expectedBody: task,
			requestBody:  `{"name": "study golang unit testing", "completed": false}`,
		},
		{
			caseName:           "error creating a task",
			expectedError:      model.ErrExecuteQuery,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().Create(task).DoAndReturn(func(newTask model.Task) (model.Task, error) {
					return task, model.ErrExecuteQuery
				}).Times(1)
			},
			requestBody: `{"name": "study golang unit testing", "completed": false}`,
		},
		{
			caseName:           "internal server error - failed preparing statement",
			expectedError:      model.ErrPreparingStatemant,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().Create(task).DoAndReturn(func(newTask model.Task) (model.Task, error) {
					return task, model.ErrPreparingStatemant
				}).Times(1)
			},
			requestBody: `{"name": "study golang unit testing", "completed": false}`,
		},
		{
			caseName:           "internal server error - failed connect to database",
			expectedError:      model.ErrConnectDatabase,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().Create(task).DoAndReturn(func(newTask model.Task) (model.Task, error) {
					return task, model.ErrConnectDatabase
				}).Times(1)
			},
			requestBody: `{"name": "study golang unit testing", "completed": false}`,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := repository.NewTaskMock(ctrl)
			ctrlMock := controller.NewTask(repoMock)

			testCase.expectedBehavior(repoMock)

			newReader := strings.NewReader(testCase.requestBody)
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/tasks", newReader)
			r.Header.Set("Authorization", pkg.GetBearerToken())

			ctrlMock.ServeHTTP(w, r)

			if testCase.expectedError != nil {
				errorMessage, _ := ioutil.ReadAll(w.Body)
				assertError(t, string(errorMessage), testCase.expectedError.Error())
				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
			} else {
				var got model.Task
				json.NewDecoder(w.Body).Decode(&got)

				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
				assertResponseBody(t, got, testCase.expectedBody)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	task := model.Task{ID: 1, Name: "study golang unit testing", Completed: true}
	cases := []struct {
		caseName           string
		expectedError      error
		expectedStatusCode int
		expectedBehavior   func(m *repository.TaskMock)
		expectedBody       model.Task
		requestBody        string
	}{
		{
			caseName:           "successfully updating a task",
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Update(task).DoAndReturn(func(modifiedTask model.Task) (model.Task, error) {
					return task, nil
				}).Times(1)
			},
			expectedBody: task,
			requestBody:  `{"id": 1, "name": "study golang unit testing", "completed": true}`,
		},
		{
			caseName:           "error updating a task",
			expectedError:      model.ErrExecuteQuery,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Update(task).DoAndReturn(func(modifiedTask model.Task) (model.Task, error) {
					return task, model.ErrExecuteQuery
				}).Times(1)
			},
			requestBody: `{"id": 1, "name": "study golang unit testing", "completed": true}`,
		},
		{
			caseName:           "internal server error - failed preparing statement",
			expectedError:      model.ErrPreparingStatemant,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Update(task).DoAndReturn(func(modifiedTask model.Task) (model.Task, error) {
					return task, model.ErrPreparingStatemant
				}).Times(1)
			},
			requestBody: `{"id": 1, "name": "study golang unit testing", "completed": true}`,
		},
		{
			caseName:           "internal server error - failed connect to database",
			expectedError:      model.ErrConnectDatabase,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Update(task).DoAndReturn(func(modifiedTask model.Task) (model.Task, error) {
					return task, model.ErrConnectDatabase
				}).Times(1)
			},
			requestBody: `{"id": 1, "name": "study golang unit testing", "completed": true}`,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := repository.NewTaskMock(ctrl)
			ctrlMock := controller.NewTask(repoMock)

			testCase.expectedBehavior(repoMock)

			newReader := strings.NewReader(testCase.requestBody)
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.ID), newReader)
			r.Header.Set("Authorization", pkg.GetBearerToken())

			ctrlMock.ServeHTTP(w, r)

			if testCase.expectedError != nil {
				errorMessage, _ := ioutil.ReadAll(w.Body)
				assertError(t, string(errorMessage), testCase.expectedError.Error())
				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
			} else {
				var got model.Task
				json.NewDecoder(w.Body).Decode(&got)

				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
				assertResponseBody(t, got, testCase.expectedBody)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	task := model.Task{ID: 1, Name: "task mock", Completed: false}
	cases := []struct {
		caseName           string
		expectedError      error
		expectedStatusCode int
		expectedBehavior   func(m *repository.TaskMock)
	}{
		{
			caseName:           "success deleteting a task",
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Delete(task.ID).DoAndReturn(func(id int) error {
					return nil
				}).Times(1)
			},
		},
		{
			caseName:           "error deleteting a task",
			expectedError:      model.ErrExecuteQuery,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Delete(task.ID).DoAndReturn(func(id int) error {
					return model.ErrExecuteQuery
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed preparing statement",
			expectedError:      model.ErrPreparingStatemant,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Delete(task.ID).DoAndReturn(func(id int) error {
					return model.ErrPreparingStatemant
				}).Times(1)
			},
		},
		{
			caseName:           "internal server error - failed connect to database",
			expectedError:      model.ErrConnectDatabase,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBehavior: func(m *repository.TaskMock) {
				m.EXPECT().FindByID(task.ID).DoAndReturn(func(id int) (model.Task, error) {
					return task, nil
				}).Times(1)
				m.EXPECT().Delete(task.ID).DoAndReturn(func(id int) error {
					return model.ErrConnectDatabase
				}).Times(1)
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := repository.NewTaskMock(ctrl)
			ctrlMock := controller.NewTask(repoMock)

			testCase.expectedBehavior(repoMock)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", task.ID), nil)
			r.Header.Set("Authorization", pkg.GetBearerToken())

			ctrlMock.ServeHTTP(w, r)

			if testCase.expectedError != nil {
				errorMessage, _ := ioutil.ReadAll(w.Body)
				assertError(t, string(errorMessage), testCase.expectedError.Error())
				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
			} else {
				assertStatusCode(t, w.Result().StatusCode, testCase.expectedStatusCode)
				assertResponseBody(t, nil, nil)
			}
		})
	}
}

func assertError(t testing.TB, got, want string) {
	t.Helper()

	if got != "" && want != "" && want != got {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
