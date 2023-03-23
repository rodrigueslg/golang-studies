package repository

import (
	"gochallenges/internal/model"
)

const (
	DbVanilla, DbOrm = "vanilla", "orm"
)

type Task interface {
	Create(task model.Task) (model.Task, error)
	FindByID(id int) (model.Task, error)
	FindByStatus(completed bool) ([]model.Task, error)
	FindAll() ([]model.Task, error)
	Update(task model.Task) (model.Task, error)
	Delete(id int) error
	Close()
}
