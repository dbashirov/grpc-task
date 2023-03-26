package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/dbashirov/grpc-tasks/api"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultPort = "8080"

func main() {

	log.Println("[INFO] start gRPC client")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("[ERROR] error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	con, err := grpc.Dial("localhost:"+port, opts)
	if err != nil {
		log.Fatalf("[ERROR] could not connect: %v", err)
	}
	defer con.Close()

	sc := api.NewTaskServiceClient(con)

	// create task
	log.Println("[INFO] creating task")
	task := &api.Task{
		Name: "Task 1",
		Desc: "More desc task 1",
	}
	createTaskRes, err := sc.CreateTask(context.Background(), &api.CreateTaskRequest{Task: task})
	if err != nil {
		log.Fatalf("[ERROR] unexpected error: %v", err)
	}
	log.Printf("[INFO] task has been created: %v", createTaskRes)
	taskID := createTaskRes.GetTask().GetId()

	// read task
	log.Println("[INFO] read task")
	readTaskReq := &api.ReadTaskRequest{Id: taskID}
	readTaskRes, err := sc.ReadTask(context.Background(), readTaskReq)
	if err != nil {
		log.Printf("[ERROR] error happened while reading: %v\n", err)
	}
	log.Printf("[INFO] task was read:%v\n", readTaskRes)

	// update task
	log.Println("[INFO] update task")
	newTask := &api.Task{
		Id:   taskID,
		Name: "Task 1",
		Desc: "More more task 1, done",
		Done: true,
	}
	updateTaskRes, err := sc.UpdateTask(context.Background(), &api.UpdateTaskRequest{Task: newTask})
	if err != nil {
		log.Printf("[ERROR] error updating task: %v\n", err)
	}
	log.Printf("[INFO] task was update: %v\n", updateTaskRes)

	// delete task
	// deleteTaskRes, err := sc.DeleteTask(context.Background(), &api.DeleteTaskRequest{Id: taskID})
	// if err != nil {
	// 	log.Printf("[ERROR] error happened while deleting: %v\n", err)
	// }
	// log.Printf("[INFO] task was deleted: %v\n", deleteTaskRes)

	// list tasks
	stream, err := sc.ListTask(context.Background(), &api.ListTaskRequest{})
	if err != nil {
		log.Fatalf("[ERROR] error while calling ListTask: %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[ERROR] something happened: %v\n", err)
		}
		log.Printf("[INFO] stream: %v\n", res.GetTask())
	}

}
