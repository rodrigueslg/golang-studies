package service

import (
	"gochallenges/internal/model"
	"gochallenges/internal/repository"
)

type Task struct {
	taskRepository repository.Task
}

func NewTask(taskRepository repository.Task) Task {
	return Task{taskRepository: taskRepository}
}

func (s *Task) Create(task model.Task) (model.Task, error) {
	if task.Name == "" {
		return task, model.ErrInvalidTaskName
	}

	if task.ID > 0 {
		existing, err := s.taskRepository.FindByID(task.ID)
		if err != nil {
			return task, err
		}
		if existing.ID > 0 {
			return task, model.ErrTaskAlreadyExists
		}
	}

	createdTask, err := s.taskRepository.Create(task)
	if err != nil {
		return task, err
	}

	return createdTask, nil
}

func (s *Task) Update(task model.Task) (model.Task, error) {
	if task.Name == "" {
		return task, model.ErrInvalidTaskName
	}

	if task.ID == 0 {
		return task, model.ErrInvalidTaskId
	}

	existing, err := s.taskRepository.FindByID(task.ID)
	if err != nil {
		return task, err
	}
	if existing.ID == 0 {
		return task, model.ErrTaskNotFound
	}

	updatedTask, err := s.taskRepository.Update(task)
	if err != nil {
		return updatedTask, err
	}

	return updatedTask, nil
}

func (s *Task) Delete(id int) error {
	if id == 0 {
		return model.ErrInvalidTaskId
	}

	existing, err := s.taskRepository.FindByID(id)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return model.ErrTaskNotFound
	}

	if err = s.taskRepository.Delete(id); err != nil {
		return err
	}

	return nil
}
