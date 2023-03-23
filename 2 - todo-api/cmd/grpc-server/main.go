package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gochallenges/api/proto"
	"gochallenges/internal/repository"
	"gochallenges/internal/service"
	"gochallenges/pkg"
)

var addr string = "0.0.0.0:5001"
var addrGw string = "0.0.0.0:5002"

func main() {
	startRpcServer()
	startRpcGateway()
}

func startRpcServer() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on %s", addr)

	s := grpc.NewServer()

	var repository = loadTaskRepository()
	pb.RegisterTasksServiceServer(s, &RpcServer{
		taskRepository: repository,
		taskService:    service.NewTask(repository),
	})

	go func() {
		log.Fatalln(s.Serve(lis))
	}()
}

func startRpcGateway() {
	conn, err := grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("failed to dial grpc server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterTasksServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gRPC-Gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    addrGw,
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on: %s", addrGw)
	log.Fatalln(gwServer.ListenAndServe())
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
