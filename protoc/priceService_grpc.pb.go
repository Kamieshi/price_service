// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: priceService.proto

package protoc

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

// OwnPriceStreamClient is the client API for OwnPriceStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OwnPriceStreamClient interface {
	GetPriceStream(ctx context.Context, in *GetPriceStreamRequest, opts ...grpc.CallOption) (OwnPriceStream_GetPriceStreamClient, error)
}

type ownPriceStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewOwnPriceStreamClient(cc grpc.ClientConnInterface) OwnPriceStreamClient {
	return &ownPriceStreamClient{cc}
}

func (c *ownPriceStreamClient) GetPriceStream(ctx context.Context, in *GetPriceStreamRequest, opts ...grpc.CallOption) (OwnPriceStream_GetPriceStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &OwnPriceStream_ServiceDesc.Streams[0], "/priceService.OwnPriceStream/GetPriceStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &ownPriceStreamGetPriceStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OwnPriceStream_GetPriceStreamClient interface {
	Recv() (*GetPriceStreamResponse, error)
	grpc.ClientStream
}

type ownPriceStreamGetPriceStreamClient struct {
	grpc.ClientStream
}

func (x *ownPriceStreamGetPriceStreamClient) Recv() (*GetPriceStreamResponse, error) {
	m := new(GetPriceStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OwnPriceStreamServer is the server API for OwnPriceStream service.
// All implementations must embed UnimplementedOwnPriceStreamServer
// for forward compatibility
type OwnPriceStreamServer interface {
	GetPriceStream(*GetPriceStreamRequest, OwnPriceStream_GetPriceStreamServer) error
	mustEmbedUnimplementedOwnPriceStreamServer()
}

// UnimplementedOwnPriceStreamServer must be embedded to have forward compatible implementations.
type UnimplementedOwnPriceStreamServer struct {
}

func (UnimplementedOwnPriceStreamServer) GetPriceStream(*GetPriceStreamRequest, OwnPriceStream_GetPriceStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetPriceStream not implemented")
}
func (UnimplementedOwnPriceStreamServer) mustEmbedUnimplementedOwnPriceStreamServer() {}

// UnsafeOwnPriceStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OwnPriceStreamServer will
// result in compilation errors.
type UnsafeOwnPriceStreamServer interface {
	mustEmbedUnimplementedOwnPriceStreamServer()
}

func RegisterOwnPriceStreamServer(s grpc.ServiceRegistrar, srv OwnPriceStreamServer) {
	s.RegisterService(&OwnPriceStream_ServiceDesc, srv)
}

func _OwnPriceStream_GetPriceStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetPriceStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OwnPriceStreamServer).GetPriceStream(m, &ownPriceStreamGetPriceStreamServer{stream})
}

type OwnPriceStream_GetPriceStreamServer interface {
	Send(*GetPriceStreamResponse) error
	grpc.ServerStream
}

type ownPriceStreamGetPriceStreamServer struct {
	grpc.ServerStream
}

func (x *ownPriceStreamGetPriceStreamServer) Send(m *GetPriceStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// OwnPriceStream_ServiceDesc is the grpc.ServiceDesc for OwnPriceStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OwnPriceStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "priceService.OwnPriceStream",
	HandlerType: (*OwnPriceStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetPriceStream",
			Handler:       _OwnPriceStream_GetPriceStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "priceService.proto",
}

// CommonPriceStreamClient is the client API for CommonPriceStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommonPriceStreamClient interface {
	GetPriceStream(ctx context.Context, in *GetPriceStreamRequest, opts ...grpc.CallOption) (CommonPriceStream_GetPriceStreamClient, error)
}

type commonPriceStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewCommonPriceStreamClient(cc grpc.ClientConnInterface) CommonPriceStreamClient {
	return &commonPriceStreamClient{cc}
}

func (c *commonPriceStreamClient) GetPriceStream(ctx context.Context, in *GetPriceStreamRequest, opts ...grpc.CallOption) (CommonPriceStream_GetPriceStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &CommonPriceStream_ServiceDesc.Streams[0], "/priceService.CommonPriceStream/GetPriceStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &commonPriceStreamGetPriceStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CommonPriceStream_GetPriceStreamClient interface {
	Recv() (*GetPriceStreamResponse, error)
	grpc.ClientStream
}

type commonPriceStreamGetPriceStreamClient struct {
	grpc.ClientStream
}

func (x *commonPriceStreamGetPriceStreamClient) Recv() (*GetPriceStreamResponse, error) {
	m := new(GetPriceStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CommonPriceStreamServer is the server API for CommonPriceStream service.
// All implementations must embed UnimplementedCommonPriceStreamServer
// for forward compatibility
type CommonPriceStreamServer interface {
	GetPriceStream(*GetPriceStreamRequest, CommonPriceStream_GetPriceStreamServer) error
	mustEmbedUnimplementedCommonPriceStreamServer()
}

// UnimplementedCommonPriceStreamServer must be embedded to have forward compatible implementations.
type UnimplementedCommonPriceStreamServer struct {
}

func (UnimplementedCommonPriceStreamServer) GetPriceStream(*GetPriceStreamRequest, CommonPriceStream_GetPriceStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetPriceStream not implemented")
}
func (UnimplementedCommonPriceStreamServer) mustEmbedUnimplementedCommonPriceStreamServer() {}

// UnsafeCommonPriceStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommonPriceStreamServer will
// result in compilation errors.
type UnsafeCommonPriceStreamServer interface {
	mustEmbedUnimplementedCommonPriceStreamServer()
}

func RegisterCommonPriceStreamServer(s grpc.ServiceRegistrar, srv CommonPriceStreamServer) {
	s.RegisterService(&CommonPriceStream_ServiceDesc, srv)
}

func _CommonPriceStream_GetPriceStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetPriceStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CommonPriceStreamServer).GetPriceStream(m, &commonPriceStreamGetPriceStreamServer{stream})
}

type CommonPriceStream_GetPriceStreamServer interface {
	Send(*GetPriceStreamResponse) error
	grpc.ServerStream
}

type commonPriceStreamGetPriceStreamServer struct {
	grpc.ServerStream
}

func (x *commonPriceStreamGetPriceStreamServer) Send(m *GetPriceStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// CommonPriceStream_ServiceDesc is the grpc.ServiceDesc for CommonPriceStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommonPriceStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "priceService.CommonPriceStream",
	HandlerType: (*CommonPriceStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetPriceStream",
			Handler:       _CommonPriceStream_GetPriceStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "priceService.proto",
}
