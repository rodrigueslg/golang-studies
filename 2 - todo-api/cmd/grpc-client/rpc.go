package main

import (
	"context"
	"log"

	pb "gochallenges/api/proto"
)

func doGetAllTasks(c pb.TasksServiceClient) {
	req := &pb.GetTasksRequest{}
	res, err := c.GetTasks(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling Task RPC: %v", err)
	}
	log.Printf("Response from Tasks RPC: %v", res)
}
