// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0--rc1
// source: auth.proto

package proto

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
	AuthHandl_IsAuthenticated_FullMethodName = "/auth.AuthHandl/IsAuthenticated"
)

// AuthHandlClient is the client API for AuthHandl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthHandlClient interface {
	IsAuthenticated(ctx context.Context, in *IsAuthRequest, opts ...grpc.CallOption) (*IsAuthResponse, error)
}

type authHandlClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthHandlClient(cc grpc.ClientConnInterface) AuthHandlClient {
	return &authHandlClient{cc}
}

func (c *authHandlClient) IsAuthenticated(ctx context.Context, in *IsAuthRequest, opts ...grpc.CallOption) (*IsAuthResponse, error) {
	out := new(IsAuthResponse)
	err := c.cc.Invoke(ctx, AuthHandl_IsAuthenticated_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthHandlServer is the server API for AuthHandl service.
// All implementations must embed UnimplementedAuthHandlServer
// for forward compatibility
type AuthHandlServer interface {
	IsAuthenticated(context.Context, *IsAuthRequest) (*IsAuthResponse, error)
	mustEmbedUnimplementedAuthHandlServer()
}

// UnimplementedAuthHandlServer must be embedded to have forward compatible implementations.
type UnimplementedAuthHandlServer struct {
}

func (UnimplementedAuthHandlServer) IsAuthenticated(context.Context, *IsAuthRequest) (*IsAuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAuthenticated not implemented")
}
func (UnimplementedAuthHandlServer) mustEmbedUnimplementedAuthHandlServer() {}

// UnsafeAuthHandlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthHandlServer will
// result in compilation errors.
type UnsafeAuthHandlServer interface {
	mustEmbedUnimplementedAuthHandlServer()
}

func RegisterAuthHandlServer(s grpc.ServiceRegistrar, srv AuthHandlServer) {
	s.RegisterService(&AuthHandl_ServiceDesc, srv)
}

func _AuthHandl_IsAuthenticated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthHandlServer).IsAuthenticated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthHandl_IsAuthenticated_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthHandlServer).IsAuthenticated(ctx, req.(*IsAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthHandl_ServiceDesc is the grpc.ServiceDesc for AuthHandl service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthHandl_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthHandl",
	HandlerType: (*AuthHandlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsAuthenticated",
			Handler:    _AuthHandl_IsAuthenticated_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
