package api

import (
	"gochallenges/internal/controller"
	"gochallenges/internal/repository"
	"net/http"
)

type HttpServer struct {
	http.Handler
	tasksController controller.Task
}

func NewServer(tasksRepository repository.Task) *HttpServer {
	s := new(HttpServer)
	s.tasksController = controller.NewTask(tasksRepository)
	s.Handler = http.HandlerFunc(s.tasksController.ServeHTTP)
	return s
}
