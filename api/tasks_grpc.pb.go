// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: api/tasks.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TaskService_CreateTask_FullMethodName = "/api.TaskService/CreateTask"
	TaskService_ReadTask_FullMethodName   = "/api.TaskService/ReadTask"
	TaskService_UpdateTask_FullMethodName = "/api.TaskService/UpdateTask"
	TaskService_DeleteTask_FullMethodName = "/api.TaskService/DeleteTask"
	TaskService_ListTask_FullMethodName   = "/api.TaskService/ListTask"
)

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponse, error)
	ReadTask(ctx context.Context, in *ReadTaskRequest, opts ...grpc.CallOption) (*ReadTaskResponse, error)
	UpdateTask(ctx context.Context, in *UpdateTaskRequest, opts ...grpc.CallOption) (*UpdateTaskResponse, error)
	DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...grpc.CallOption) (*DeleteTaskResponse, error)
	ListTask(ctx context.Context, in *ListTaskRequest, opts ...grpc.CallOption) (TaskService_ListTaskClient, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponse, error) {
	out := new(CreateTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_CreateTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) ReadTask(ctx context.Context, in *ReadTaskRequest, opts ...grpc.CallOption) (*ReadTaskResponse, error) {
	out := new(ReadTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_ReadTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) UpdateTask(ctx context.Context, in *UpdateTaskRequest, opts ...grpc.CallOption) (*UpdateTaskResponse, error) {
	out := new(UpdateTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_UpdateTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...grpc.CallOption) (*DeleteTaskResponse, error) {
	out := new(DeleteTaskResponse)
	err := c.cc.Invoke(ctx, TaskService_DeleteTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) ListTask(ctx context.Context, in *ListTaskRequest, opts ...grpc.CallOption) (TaskService_ListTaskClient, error) {
	stream, err := c.cc.NewStream(ctx, &TaskService_ServiceDesc.Streams[0], TaskService_ListTask_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &taskServiceListTaskClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TaskService_ListTaskClient interface {
	Recv() (*ListTaskResponse, error)
	grpc.ClientStream
}

type taskServiceListTaskClient struct {
	grpc.ClientStream
}

func (x *taskServiceListTaskClient) Recv() (*ListTaskResponse, error) {
	m := new(ListTaskResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations should embed UnimplementedTaskServiceServer
// for forward compatibility
type TaskServiceServer interface {
	CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponse, error)
	ReadTask(context.Context, *ReadTaskRequest) (*ReadTaskResponse, error)
	UpdateTask(context.Context, *UpdateTaskRequest) (*UpdateTaskResponse, error)
	DeleteTask(context.Context, *DeleteTaskRequest) (*DeleteTaskResponse, error)
	ListTask(*ListTaskRequest, TaskService_ListTaskServer) error
}

// UnimplementedTaskServiceServer should be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (UnimplementedTaskServiceServer) CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}
func (UnimplementedTaskServiceServer) ReadTask(context.Context, *ReadTaskRequest) (*ReadTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadTask not implemented")
}
func (UnimplementedTaskServiceServer) UpdateTask(context.Context, *UpdateTaskRequest) (*UpdateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTask not implemented")
}
func (UnimplementedTaskServiceServer) DeleteTask(context.Context, *DeleteTaskRequest) (*DeleteTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTask not implemented")
}
func (UnimplementedTaskServiceServer) ListTask(*ListTaskRequest, TaskService_ListTaskServer) error {
	return status.Errorf(codes.Unimplemented, "method ListTask not implemented")
}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_CreateTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).CreateTask(ctx, req.(*CreateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_ReadTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).ReadTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_ReadTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).ReadTask(ctx, req.(*ReadTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_UpdateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).UpdateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_UpdateTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).UpdateTask(ctx, req.(*UpdateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_DeleteTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).DeleteTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_DeleteTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).DeleteTask(ctx, req.(*DeleteTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_ListTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskServiceServer).ListTask(m, &taskServiceListTaskServer{stream})
}

type TaskService_ListTaskServer interface {
	Send(*ListTaskResponse) error
	grpc.ServerStream
}

type taskServiceListTaskServer struct {
	grpc.ServerStream
}

func (x *taskServiceListTaskServer) Send(m *ListTaskResponse) error {
	return x.ServerStream.SendMsg(m)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTask",
			Handler:    _TaskService_CreateTask_Handler,
		},
		{
			MethodName: "ReadTask",
			Handler:    _TaskService_ReadTask_Handler,
		},
		{
			MethodName: "UpdateTask",
			Handler:    _TaskService_UpdateTask_Handler,
		},
		{
			MethodName: "DeleteTask",
			Handler:    _TaskService_DeleteTask_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListTask",
			Handler:       _TaskService_ListTask_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/tasks.proto",
}
