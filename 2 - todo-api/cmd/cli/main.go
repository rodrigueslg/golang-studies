package main

import (
	repo "gochallenges/internal/repository"
	"gochallenges/pkg"
	"log"
)

func main() {
	dbConfig := pkg.GetDbConfig()
	dbConnn := pkg.GetMysqlDbConnection(dbConfig)

	tasksRepo, err := repo.NewTaskSql(dbConfig.Driver, dbConnn)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	tasks, err := tasksRepo.FindAll()
	if err != nil {
		log.Fatalf("Could not get tasks: %s", err)
	}

	for _, task := range tasks {
		log.Printf("%+v", task)
	}
}
