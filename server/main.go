package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "github.com/dbashirov/grpc-tasks/api"
	"github.com/joho/godotenv"
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
	id   primitive.ObjectID `bson:"_id,omitempty"`
	name string             `bson:"name"`
	desc string             `bsoon:"desc"`
	done bool               `bson:"done"`
}

type server struct {
	pb.TaskServiceServer
}

func (s *server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {

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

	return &pb.CreateTaskResponse{
		Task: &pb.Task{
			Id:   oid.Hex(),
			Name: t.GetName(),
			Desc: t.GetDesc(),
			Done: t.GetDone(),
		},
	}, nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load(".env")
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
	// log.Println(mongoURL)
	// log.Println(options.Client().ApplyURI(mongoURL))

	collection = (*mongo.Collection)(client.Database("taskdb").Collection("task"))

	log.Println("Task service started")
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{})

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

	// reflection.Register(s)
}
