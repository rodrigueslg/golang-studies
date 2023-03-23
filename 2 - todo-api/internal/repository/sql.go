package repository

import (
	"database/sql"
	"gochallenges/internal/model"

	_ "github.com/go-sql-driver/mysql"
)

type TaskSql struct {
	DB *sql.DB
}

func NewTaskSql(driverName string, connStr string) (Task, error) {
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, model.ErrConnectDatabase
	}

	err = db.Ping()
	if err != nil {
		return nil, model.ErrConnectDatabase
	}

	return &TaskSql{db}, nil
}

func (r *TaskSql) Create(task model.Task) (model.Task, error) {
	statement, err := r.DB.Prepare("INSERT INTO task (name, completed) VALUES (?, ?)")
	if err != nil {
		return task, model.ErrPreparingStatemant
	}
	defer statement.Close()

	result, err := statement.Exec(task.Name, task.Completed)
	if err != nil {
		return task, model.ErrExecuteQuery
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return task, model.ErrExecuteQuery
	}

	task.ID = int(newId)
	return task, nil
}

func (r *TaskSql) FindByID(taskId int) (model.Task, error) {
	var task model.Task

	row, err := r.DB.Query("SELECT id, name, completed FROM task WHERE id = ?", taskId)
	if err != nil {
		return task, model.ErrExecuteQuery
	}
	defer row.Close()

	if row.Next() {
		if err := row.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return task, model.ErrScanningRows
		}
	}

	return task, nil
}

func (r *TaskSql) FindByStatus(completed bool) ([]model.Task, error) {
	rows, err := r.DB.Query("SELECT id, name, completed FROM task WHERE completed = ?", completed)
	if err != nil {
		return nil, model.ErrExecuteQuery
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, model.ErrScanningRows
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskSql) FindAll() ([]model.Task, error) {
	rows, err := r.DB.Query("SELECT id, name, completed FROM task")
	if err != nil {
		return nil, model.ErrExecuteQuery
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, model.ErrScanningRows
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskSql) Update(task model.Task) (model.Task, error) {
	var updatedTask model.Task

	statement, err := r.DB.Prepare("UPDATE task SET name = ?, completed = ? WHERE id = ?")
	if err != nil {
		return updatedTask, model.ErrPreparingStatemant
	}
	defer statement.Close()

	if _, err := statement.Exec(task.Name, task.Completed, task.ID); err != nil {
		return updatedTask, model.ErrExecuteQuery
	}

	updatedTask, err = r.FindByID(task.ID)
	if err != nil {
		return updatedTask, err
	}

	return updatedTask, nil
}

func (r *TaskSql) Delete(id int) error {
	statement, err := r.DB.Prepare("DELETE FROM task WHERE id = ?")
	if err != nil {
		return model.ErrPreparingStatemant
	}
	defer statement.Close()

	if _, err := statement.Exec(id); err != nil {
		return model.ErrExecuteQuery
	}

	return nil
}

func (r *TaskSql) Close() {
	r.DB.Close()
}
