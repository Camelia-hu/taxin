// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.4
// source: system.proto

package systempb

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
	SystemService_SendFile_FullMethodName = "/system.SystemService/SendFile"
)

// SystemServiceClient is the client API for SystemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SystemServiceClient interface {
	// 读取文件以流形式返回
	SendFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (SystemService_SendFileClient, error)
}

type systemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSystemServiceClient(cc grpc.ClientConnInterface) SystemServiceClient {
	return &systemServiceClient{cc}
}

func (c *systemServiceClient) SendFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (SystemService_SendFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &SystemService_ServiceDesc.Streams[0], SystemService_SendFile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &systemServiceSendFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SystemService_SendFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type systemServiceSendFileClient struct {
	grpc.ClientStream
}

func (x *systemServiceSendFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SystemServiceServer is the server API for SystemService service.
// All implementations must embed UnimplementedSystemServiceServer
// for forward compatibility
type SystemServiceServer interface {
	// 读取文件以流形式返回
	SendFile(*FileRequest, SystemService_SendFileServer) error
	mustEmbedUnimplementedSystemServiceServer()
}

// UnimplementedSystemServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSystemServiceServer struct {
}

func (UnimplementedSystemServiceServer) SendFile(*FileRequest, SystemService_SendFileServer) error {
	return status.Errorf(codes.Unimplemented, "method SendFile not implemented")
}
func (UnimplementedSystemServiceServer) mustEmbedUnimplementedSystemServiceServer() {}

// UnsafeSystemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SystemServiceServer will
// result in compilation errors.
type UnsafeSystemServiceServer interface {
	mustEmbedUnimplementedSystemServiceServer()
}

func RegisterSystemServiceServer(s grpc.ServiceRegistrar, srv SystemServiceServer) {
	s.RegisterService(&SystemService_ServiceDesc, srv)
}

func _SystemService_SendFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SystemServiceServer).SendFile(m, &systemServiceSendFileServer{stream})
}

type SystemService_SendFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type systemServiceSendFileServer struct {
	grpc.ServerStream
}

func (x *systemServiceSendFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

// SystemService_ServiceDesc is the grpc.ServiceDesc for SystemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SystemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "system.SystemService",
	HandlerType: (*SystemServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendFile",
			Handler:       _SystemService_SendFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "system.proto",
}
