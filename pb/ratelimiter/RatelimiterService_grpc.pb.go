// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: pb/RatelimiterService.proto

package ratelimiter

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Ratelimiter_Allow_FullMethodName = "/ratelimiter.Ratelimiter/Allow"
	Ratelimiter_Clear_FullMethodName = "/ratelimiter.Ratelimiter/Clear"
)

// RatelimiterClient is the client API for Ratelimiter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RatelimiterClient interface {
	Allow(ctx context.Context, in *AllowRequest, opts ...grpc.CallOption) (*AllowResponse, error)
	Clear(ctx context.Context, in *ClearRequest, opts ...grpc.CallOption) (*Empty, error)
}

type ratelimiterClient struct {
	cc grpc.ClientConnInterface
}

func NewRatelimiterClient(cc grpc.ClientConnInterface) RatelimiterClient {
	return &ratelimiterClient{cc}
}

func (c *ratelimiterClient) Allow(ctx context.Context, in *AllowRequest, opts ...grpc.CallOption) (*AllowResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllowResponse)
	err := c.cc.Invoke(ctx, Ratelimiter_Allow_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratelimiterClient) Clear(ctx context.Context, in *ClearRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, Ratelimiter_Clear_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RatelimiterServer is the server API for Ratelimiter service.
// All implementations must embed UnimplementedRatelimiterServer
// for forward compatibility.
type RatelimiterServer interface {
	Allow(context.Context, *AllowRequest) (*AllowResponse, error)
	Clear(context.Context, *ClearRequest) (*Empty, error)
	mustEmbedUnimplementedRatelimiterServer()
}

// UnimplementedRatelimiterServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRatelimiterServer struct{}

func (UnimplementedRatelimiterServer) Allow(context.Context, *AllowRequest) (*AllowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allow not implemented")
}
func (UnimplementedRatelimiterServer) Clear(context.Context, *ClearRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Clear not implemented")
}
func (UnimplementedRatelimiterServer) mustEmbedUnimplementedRatelimiterServer() {}
func (UnimplementedRatelimiterServer) testEmbeddedByValue()                     {}

// UnsafeRatelimiterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RatelimiterServer will
// result in compilation errors.
type UnsafeRatelimiterServer interface {
	mustEmbedUnimplementedRatelimiterServer()
}

func RegisterRatelimiterServer(s grpc.ServiceRegistrar, srv RatelimiterServer) {
	// If the following call pancis, it indicates UnimplementedRatelimiterServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Ratelimiter_ServiceDesc, srv)
}

func _Ratelimiter_Allow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatelimiterServer).Allow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ratelimiter_Allow_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatelimiterServer).Allow(ctx, req.(*AllowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ratelimiter_Clear_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatelimiterServer).Clear(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ratelimiter_Clear_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatelimiterServer).Clear(ctx, req.(*ClearRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Ratelimiter_ServiceDesc is the grpc.ServiceDesc for Ratelimiter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ratelimiter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ratelimiter.Ratelimiter",
	HandlerType: (*RatelimiterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Allow",
			Handler:    _Ratelimiter_Allow_Handler,
		},
		{
			MethodName: "Clear",
			Handler:    _Ratelimiter_Clear_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/RatelimiterService.proto",
}
