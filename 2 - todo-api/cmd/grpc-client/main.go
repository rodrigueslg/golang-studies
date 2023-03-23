package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gochallenges/api/proto"
)

var addr string = "localhost:5001"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failded to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTasksServiceClient(conn)
	doGetAllTasks(c)
}
