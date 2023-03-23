package repository

import (
	"gochallenges/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TaskOrm struct {
	db *gorm.DB
}

func NewTaskOrm(connStr string) (Task, error) {
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(model.ErrConnectDatabase)
	}

	db.AutoMigrate(&TaskOrm{})

	return &TaskOrm{db}, nil
}

func (r *TaskOrm) Create(task model.Task) (model.Task, error) {
	newTask := model.Task{ID: task.ID, Name: task.Name, Completed: task.Completed}
	if err := r.db.Create(&newTask).Error; err != nil {
		return newTask, model.ErrInsertingRow
	}

	return newTask, nil
}

func (r *TaskOrm) FindByID(id int) (model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, id).Error; err != nil {
		return task, model.ErrExecuteQuery
	}

	return task, nil
}

func (r *TaskOrm) FindByStatus(completed bool) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.Where("Completed = ?", completed).Find(&tasks).Error; err != nil {
		return nil, model.ErrExecuteQuery
	}

	return tasks, nil
}

func (r *TaskOrm) FindAll() ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, model.ErrExecuteQuery
	}

	return tasks, nil
}

func (r *TaskOrm) Update(task model.Task) (model.Task, error) {
	var taskOrm model.Task
	if err := r.db.First(&taskOrm, task.ID).Error; err != nil {
		return taskOrm, model.ErrExecuteQuery
	}

	taskOrm.Name = task.Name
	taskOrm.Completed = task.Completed
	r.db.Save(&taskOrm)

	return taskOrm, nil
}

func (r *TaskOrm) Delete(id int) error {
	if err := r.db.Delete(&model.Task{}, id).Error; err != nil {
		return model.ErrExecuteQuery
	}
	return nil
}

func (r *TaskOrm) Close() {
}
