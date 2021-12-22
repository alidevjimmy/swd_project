// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package centerpb

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

// CenterServiceClient is the client API for CenterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CenterServiceClient interface {
	FindAll(ctx context.Context, in *FindAllRequest, opts ...grpc.CallOption) (*FindAllResponse, error)
	Find(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error)
}

type centerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCenterServiceClient(cc grpc.ClientConnInterface) CenterServiceClient {
	return &centerServiceClient{cc}
}

func (c *centerServiceClient) FindAll(ctx context.Context, in *FindAllRequest, opts ...grpc.CallOption) (*FindAllResponse, error) {
	out := new(FindAllResponse)
	err := c.cc.Invoke(ctx, "/center.CenterService/FindAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centerServiceClient) Find(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error) {
	out := new(FindResponse)
	err := c.cc.Invoke(ctx, "/center.CenterService/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CenterServiceServer is the server API for CenterService service.
// All implementations must embed UnimplementedCenterServiceServer
// for forward compatibility
type CenterServiceServer interface {
	FindAll(context.Context, *FindAllRequest) (*FindAllResponse, error)
	Find(context.Context, *FindRequest) (*FindResponse, error)
	mustEmbedUnimplementedCenterServiceServer()
}

// UnimplementedCenterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCenterServiceServer struct {
}

func (UnimplementedCenterServiceServer) FindAll(context.Context, *FindAllRequest) (*FindAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAll not implemented")
}
func (UnimplementedCenterServiceServer) Find(context.Context, *FindRequest) (*FindResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedCenterServiceServer) mustEmbedUnimplementedCenterServiceServer() {}

// UnsafeCenterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CenterServiceServer will
// result in compilation errors.
type UnsafeCenterServiceServer interface {
	mustEmbedUnimplementedCenterServiceServer()
}

func RegisterCenterServiceServer(s grpc.ServiceRegistrar, srv CenterServiceServer) {
	s.RegisterService(&CenterService_ServiceDesc, srv)
}

func _CenterService_FindAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CenterServiceServer).FindAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/center.CenterService/FindAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CenterServiceServer).FindAll(ctx, req.(*FindAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CenterService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CenterServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/center.CenterService/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CenterServiceServer).Find(ctx, req.(*FindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CenterService_ServiceDesc is the grpc.ServiceDesc for CenterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CenterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "center.CenterService",
	HandlerType: (*CenterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAll",
			Handler:    _CenterService_FindAll_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _CenterService_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/pbs/centerpb/center.proto",
}