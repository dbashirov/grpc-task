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

	// "google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	defaultPort = "8080"
	collection  *mongo.Collection
)

type task struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Desc string             `bsoon:"desc"`
	Done bool               `bson:"done"`
}

func newTask() *task {
	return &task{}
}

func getTaskGRPC(data *task) *api.Task {
	return &api.Task{
		Id:   data.ID.Hex(),
		Name: data.Name,
		Desc: data.Desc,
		Done: data.Done,
	}
}

type server struct {
	api.TaskServiceServer
}

func (*server) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {

	log.Println("[INFO] start create task")

	t := req.GetTask()
	data := task{
		Name: t.GetName(),
		Desc: t.GetDesc(),
		Done: t.GetDone(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("[ERROR] internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"[ERROR] Cannot convert to OID",
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

	log.Println("[INFO] read task")

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"[ERROR] cannot parse ID",
		)
	}

	data := newTask()
	filter := bson.M{"_id": oid}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("[ERROR] cannot find task with ID: %v", err),
		)
	}

	return &api.ReadTaskResponse{
		Task: getTaskGRPC(data),
	}, nil
}

func (*server) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.UpdateTaskResponse, error) {

	log.Println("[INFO] update task")

	t := req.GetTask()
	oid, err := primitive.ObjectIDFromHex(t.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"[ERROR] cannot parse ID",
		)
	}

	data := newTask()
	filter := bson.M{"_id": oid}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("[ERROR] cannot find task with ID: %v", err),
		)
	}

	data.Name = t.GetName()
	data.Desc = t.GetDesc()
	data.Done = t.GetDone()

	_, err = collection.ReplaceOne(ctx, filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("[ERROR] cannot update object in MongoDB: %v", err),
		)
	}

	return &api.UpdateTaskResponse{
		Task: getTaskGRPC(data),
	}, nil
}

func (*server) DeleteTask(ctx context.Context, req *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {

	log.Println("[INFO] delete task")

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"[ERROR] cannot parse ID",
		)
	}

	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("[ERROR] cannot delete object in MongoDB: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"[ERROR] cannot fint task in MongoDB",
		)
	}

	return &api.DeleteTaskResponse{
		Id: req.GetId(),
	}, nil
}

func (*server) ListTask(_ *api.ListTaskRequest, stream api.TaskService_ListTaskServer) error {

	log.Println("[INFO] stream list tasks")

	cur, err := collection.Find(context.Background(), primitive.D{})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("[ERROR] unknown internal error: %v\n", err),
		)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := newTask()
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("[ERROR] error while decoding data from MongoDB: %v\n", err),
			)
		}
		stream.Send(&api.ListTaskResponse{
			Task: getTaskGRPC(data),
		})
	}

	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("[ERROR] unknown internal error: %v", err),
		)
	}

	return nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("[ERROR] error loading .env files")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mongoURL := os.Getenv("MONGODB_URL")

	log.Println("[INFO] connect to MongoDB")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	collection = (*mongo.Collection)(client.Database("taskdb").Collection("task"))

	log.Println("[INFO] task service started")
	s := grpc.NewServer()
	api.RegisterTaskServiceServer(s, &server{})

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("[INFO] starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serv: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Println("[INFO] closing MongoDB connection")
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatalf("Error on disconnection with MongoDB: %v\n", err)
	}

	log.Println("[INFO] stoping server")
	s.Stop()
	log.Println("[INFO] end of Program")

	// TODO
	// reflection.Register(s)
}
