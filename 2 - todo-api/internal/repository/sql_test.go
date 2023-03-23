package repository_test

import (
	"database/sql"
	"gochallenges/internal/model"
	"gochallenges/internal/repository"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var taskMock = model.Task{
	ID:        1,
	Name:      "Task 1",
	Completed: false,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestCreate(t *testing.T) {
	db, mock := NewMock()
	repo := &repository.TaskSql{db}
	defer func() {
		repo.Close()
	}()

	query := "INSERT INTO task \\(name, completed\\) VALUES \\(\\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(taskMock.Name, taskMock.Completed).WillReturnResult(sqlmock.NewResult(0, 1))

	task, err := repo.Create(taskMock)
	if err != nil {
		t.Errorf("Error was not expected while creating task, got %s", err)
	}
	if task.ID != taskMock.ID {
		t.Errorf("Expected task ID to be %d, got %d", taskMock.ID, task.ID)
	}
}

func TestFindByID(t *testing.T) {
	db, mock := NewMock()
	repo := &repository.TaskSql{db}
	defer func() {
		repo.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "completed"}).AddRow(taskMock.ID, taskMock.Name, taskMock.Completed)

	query := "SELECT id, name, completed FROM task WHERE id = \\?"
	mock.ExpectQuery(query).WithArgs(taskMock.ID).WillReturnRows(rows)

	_, err := repo.FindByID(taskMock.ID)
	if err != nil {
		t.Errorf("Error was not expected while creating task, got %s", err)
	}
}
