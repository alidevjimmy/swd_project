// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package consultantpb

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

// ConsultantServiceClient is the client API for ConsultantService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsultantServiceClient interface {
	FindConsultant(ctx context.Context, in *FindConsultantRequest, opts ...grpc.CallOption) (*FindConsultantResponse, error)
	FindAllConsultants(ctx context.Context, in *FindAllConsultantsRequest, opts ...grpc.CallOption) (*FindAllConsultantsResponse, error)
}

type consultantServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConsultantServiceClient(cc grpc.ClientConnInterface) ConsultantServiceClient {
	return &consultantServiceClient{cc}
}

func (c *consultantServiceClient) FindConsultant(ctx context.Context, in *FindConsultantRequest, opts ...grpc.CallOption) (*FindConsultantResponse, error) {
	out := new(FindConsultantResponse)
	err := c.cc.Invoke(ctx, "/consultant.ConsultantService/FindConsultant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consultantServiceClient) FindAllConsultants(ctx context.Context, in *FindAllConsultantsRequest, opts ...grpc.CallOption) (*FindAllConsultantsResponse, error) {
	out := new(FindAllConsultantsResponse)
	err := c.cc.Invoke(ctx, "/consultant.ConsultantService/FindAllConsultants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsultantServiceServer is the server API for ConsultantService service.
// All implementations must embed UnimplementedConsultantServiceServer
// for forward compatibility
type ConsultantServiceServer interface {
	FindConsultant(context.Context, *FindConsultantRequest) (*FindConsultantResponse, error)
	FindAllConsultants(context.Context, *FindAllConsultantsRequest) (*FindAllConsultantsResponse, error)
	mustEmbedUnimplementedConsultantServiceServer()
}

// UnimplementedConsultantServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConsultantServiceServer struct {
}

func (UnimplementedConsultantServiceServer) FindConsultant(context.Context, *FindConsultantRequest) (*FindConsultantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindConsultant not implemented")
}
func (UnimplementedConsultantServiceServer) FindAllConsultants(context.Context, *FindAllConsultantsRequest) (*FindAllConsultantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllConsultants not implemented")
}
func (UnimplementedConsultantServiceServer) mustEmbedUnimplementedConsultantServiceServer() {}

// UnsafeConsultantServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsultantServiceServer will
// result in compilation errors.
type UnsafeConsultantServiceServer interface {
	mustEmbedUnimplementedConsultantServiceServer()
}

func RegisterConsultantServiceServer(s grpc.ServiceRegistrar, srv ConsultantServiceServer) {
	s.RegisterService(&ConsultantService_ServiceDesc, srv)
}

func _ConsultantService_FindConsultant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindConsultantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsultantServiceServer).FindConsultant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/consultant.ConsultantService/FindConsultant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsultantServiceServer).FindConsultant(ctx, req.(*FindConsultantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsultantService_FindAllConsultants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindAllConsultantsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsultantServiceServer).FindAllConsultants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/consultant.ConsultantService/FindAllConsultants",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsultantServiceServer).FindAllConsultants(ctx, req.(*FindAllConsultantsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConsultantService_ServiceDesc is the grpc.ServiceDesc for ConsultantService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConsultantService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "consultant.ConsultantService",
	HandlerType: (*ConsultantServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindConsultant",
			Handler:    _ConsultantService_FindConsultant_Handler,
		},
		{
			MethodName: "FindAllConsultants",
			Handler:    _ConsultantService_FindAllConsultants_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/pbs/consultantpb/consultant.proto",
}