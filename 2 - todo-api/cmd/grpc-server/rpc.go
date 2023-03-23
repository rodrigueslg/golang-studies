package main

import (
	"context"
	"gochallenges/internal/model"
	"gochallenges/internal/repository"
	"gochallenges/internal/service"

	"github.com/golang/protobuf/ptypes/empty"

	pb "gochallenges/api/proto"
)

type RpcServer struct {
	pb.TasksServiceServer
	taskRepository repository.Task
	taskService    service.Task
}

func (s *RpcServer) GetTasks(ctx context.Context, in *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	var tasks []model.Task
	var err error

	if in.String() == "" {
		tasks, err = s.taskRepository.FindAll()
	} else {
		tasks, err = s.taskRepository.FindByStatus(in.GetCompleted())
	}
	if err != nil {
		return nil, err
	}

	var pbTasks []*pb.Task
	for _, task := range tasks {
		pbTasks = append(pbTasks, &pb.Task{
			Id:        int32(task.ID),
			Name:      task.Name,
			Completed: task.Completed,
		})
	}

	return &pb.GetTasksResponse{
		Tasks: pbTasks,
	}, nil
}

func (s *RpcServer) GetTaskById(ctx context.Context, in *pb.GetTasksByIdRequest) (*pb.GetTasksByIdResponse, error) {
	task, err := s.taskRepository.FindByID(int(in.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetTasksByIdResponse{
		Task: &pb.Task{
			Id:        int32(task.ID),
			Name:      task.Name,
			Completed: task.Completed,
		},
	}, nil
}

func (s *RpcServer) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	var task = model.Task{
		Name:      in.GetTask().Name,
		Completed: in.GetTask().Completed,
	}

	createdTask, err := s.taskService.Create(task)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskResponse{
		Task: &pb.Task{
			Id:        int32(createdTask.ID),
			Name:      createdTask.Name,
			Completed: createdTask.Completed,
		},
	}, nil
}

func (s *RpcServer) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	var task = model.Task{
		ID:        int(in.GetTask().Id),
		Name:      in.GetTask().Name,
		Completed: in.GetTask().Completed,
	}

	updatedTask, err := s.taskService.Update(task)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTaskResponse{
		Task: &pb.Task{
			Id:        int32(updatedTask.ID),
			Name:      updatedTask.Name,
			Completed: updatedTask.Completed,
		},
	}, nil
}

func (s *RpcServer) DeleteTask(ctx context.Context, in *pb.DeleteTaskRequest) (*empty.Empty, error) {
	err := s.taskService.Delete(int(in.Id))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
