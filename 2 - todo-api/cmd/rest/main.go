package main

import (
	"gochallenges/internal/api"
	"gochallenges/internal/repository"
	"gochallenges/pkg"
	"log"
	"net/http"
)

func main() {
	server := api.NewServer(loadTaskRepository())
	err := http.ListenAndServe(":5000", server.Handler)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func loadTaskRepository() repository.Task {
	dbConfig := pkg.GetDbConfig()
	dbConnn := pkg.GetMysqlDbConnection(dbConfig)
	dbImpl := pkg.GetDbImplementation()

	var taskRepository repository.Task
	var err error

	switch dbImpl {
	case repository.DbVanilla:
		taskRepository, err = repository.NewTaskSql(dbConfig.Driver, dbConnn)
	case repository.DbOrm:
		taskRepository, err = repository.NewTaskOrm(dbConnn)
	default:
		panic("Invalid db implementation")
	}

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	return taskRepository
}
