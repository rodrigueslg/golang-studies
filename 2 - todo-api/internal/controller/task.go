package controller

import (
	"errors"
	"gochallenges/internal/model"
	"gochallenges/internal/repository"
	"gochallenges/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	repository repository.Task
	service    service.Task
}

func NewTask(repository repository.Task) Task {
	return Task{
		repository: repository,
		service:    service.NewTask(repository),
	}
}

func (c *Task) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorized := authorizeRequest(w, r)
	if !authorized {
		writeUnauthorizedResponse(w, model.ErrUnauthorized)
		return
	}

	if !strings.Contains(r.URL.Path, "/tasks") {
		writeNotFoundResponse(w)
		return
	}

	switch r.Method {

	case http.MethodGet:
		if id, length := GetTaskIdFromRequest(r.URL.Path); length > 0 && id != 0 {
			c.GetById(w, id)
		} else {
			if status, length := GetTaskStatusFromRequest(r); length > 0 {
				c.GetByStatus(w, status)
			} else {
				c.GetAll(w)
			}
		}

	case http.MethodPost:
		c.Create(w, r)

	case http.MethodPut:
		if id, length := GetTaskIdFromRequest(r.URL.Path); length > 0 {
			c.Update(w, r, id)
		}

	case http.MethodDelete:
		if id, length := GetTaskIdFromRequest(r.URL.Path); length > 0 {
			c.Delete(w, r, id)
		}
	}
}

func (c *Task) GetAll(w http.ResponseWriter) {
	tasks, err := c.repository.FindAll()
	if err != nil {
		writeInternalErrorResponse(w, err)
		return
	}
	writeOkResponse(w, tasks)
}

func (c *Task) GetById(w http.ResponseWriter, id int) {
	task, err := c.repository.FindByID(id)
	if err != nil {
		writeInternalErrorResponse(w, err)
		return
	}
	writeOkResponse(w, task)
}

func (c *Task) GetByStatus(w http.ResponseWriter, completed bool) {
	task, err := c.repository.FindByStatus(completed)
	if err != nil {
		writeInternalErrorResponse(w, err)
		return
	}
	writeOkResponse(w, task)
}

func (c *Task) Create(w http.ResponseWriter, r *http.Request) {
	task := model.Task{}
	if err := parseJsonBody(w, r, &task); err != nil {
		writeBadRequestResponse(w, model.ErrInvalidRequestBody)
		return
	}

	createdTask, err := c.service.Create(task)
	if err != nil {
		if errors.Is(err, model.ErrInvalidTaskName) || errors.Is(err, model.ErrTaskAlreadyExists) {
			writeBadRequestResponse(w, err)
			return
		}
		writeInternalErrorResponse(w, err)
		return
	}

	writeCreatedResponse(w, createdTask)
}

func (c *Task) Update(w http.ResponseWriter, r *http.Request, id int) {
	modifiedTask := model.Task{}
	if err := parseJsonBody(w, r, &modifiedTask); err != nil {
		writeBadRequestResponse(w, model.ErrInvalidRequestBody)
		return
	}

	updatedTask, err := c.service.Update(modifiedTask)
	if err != nil {
		if errors.Is(err, model.ErrInvalidTaskId) || errors.Is(err, model.ErrInvalidTaskName) || errors.Is(err, model.ErrTaskNotFound) {
			writeBadRequestResponse(w, err)
			return
		}
		writeInternalErrorResponse(w, err)
		return
	}

	writeOkResponse(w, updatedTask)
}

func (c *Task) Delete(w http.ResponseWriter, r *http.Request, taskId int) {
	err := c.service.Delete(taskId)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			writeBadRequestResponse(w, err)
			return
		}
		writeInternalErrorResponse(w, err)
		return
	}

	writeOkResponse(w, nil)
}

func GetTaskIdFromRequest(url string) (int, int) {
	idString := strings.TrimPrefix(url, "/tasks/")
	idLength := len(idString)
	id, _ := strconv.Atoi(idString)
	return id, idLength
}

func GetTaskStatusFromRequest(r *http.Request) (bool, int) {
	paramString := r.URL.Query().Get("completed")
	paramLength := len(paramString)
	status, _ := strconv.ParseBool(paramString)
	return status, paramLength
}
