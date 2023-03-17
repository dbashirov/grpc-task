package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/dbashirov/grpc-tasks/api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	defaultPort = "8080"
	collection  *mongo.Collection
)

type task struct {
	id   primitive.ObjectID `bson:"_id,omitempty"`
	name string             `bson:"name"`
	desc string             `bsoon:"desc"`
	done bool               `bson:"done"`
}

func getTaskGRPC(data *task) *api.Task {
	return &api.Task{
		Id:   data.id.Hex(),
		Name: data.name,
		Desc: data.desc,
		Done: data.done,
	}
}

type server struct {
	api.TaskServiceServer
}

func (*server) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {

	log.Println("Start create task")

	t := req.GetTask()
	data := task{
		name: t.GetName(),
		desc: t.GetDesc(),
		done: t.GetDone(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot convert to OID",
		)
	}

	log.Println("End create task")

	return &api.CreateTaskResponse{
		// Task: getTaskData(&data),
		Task: &api.Task{
			Id:   oid.Hex(),
			Name: t.GetName(),
			Desc: t.GetDesc(),
			Done: t.GetDone(),
		},
	}, nil
}

func (*server) ReadTask(ctx context.Context, req *api.ReadTaskRequest) (*api.ReadTaskResponse, error) {

	log.Println("Read task")

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	data := &task{}
	filter := bson.M{"_id": oid}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find task with ID: %v", err),
		)
	}

	return &api.ReadTaskResponse{
		Task: getTaskGRPC(data),
	}, nil
}

func (*server) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.UpdateTaskResponse, error) {

	log.Println("Update task")

	t := req.GetTask()
	oid, err := primitive.ObjectIDFromHex(t.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	data := &task{}
	filter := bson.M{"_id": oid}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find task with ID: %v", err),
		)
	}

	data.name = t.GetName()
	data.desc = t.GetDesc()
	data.done = t.GetDone()

	_, err = collection.ReplaceOne(ctx, filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", err),
		)
	}

	return &api.UpdateTaskResponse{
		Task: getTaskGRPC(data),
	}, nil
}

func (*server) DeleteTask(ctx context.Context, req *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {

	log.Println("Delete task")

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot fint task in MongoDB",
		)
	}

	return &api.DeleteTaskResponse{
		Id: req.GetId(),
	}, nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load("env_test.env")
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mongoURL := os.Getenv("MONGODB_URL")

	log.Println("Connect to MongoDB")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	collection = (*mongo.Collection)(client.Database("taskdb").Collection("task"))

	log.Println("Task service started")
	s := grpc.NewServer()
	api.RegisterTaskServiceServer(s, &server{})

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serv: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Println("Closing MongoDB connection")
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatalf("Error on disconnection with MongoDB: %v\n", err)
	}

	log.Println("Stoping server")
	s.Stop()
	log.Println("End of Program")

	reflection.Register(s)
}
