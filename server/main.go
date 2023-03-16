package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dbashirov/grpc-tasks/api"
	pb "github.com/dbashirov/grpc-tasks/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

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
			fmt.Sprintf("Connot convert to OID"),
		)
	}

	log.Println("End create task")

	return &api.CreateTaskResponse{
		Task: &api.Task{
			Id:   oid.Hex(),
			Name: t.GetName(),
			Desc: t.GetDesc(),
			Done: t.GetDone(),
		},
	}, nil
}

func main() {

}
