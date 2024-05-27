// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: warehouse.proto

package warehouse

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
	WareHouse_SetValue_FullMethodName    = "/warehouse.WareHouse/SetValue"
	WareHouse_GetValue_FullMethodName    = "/warehouse.WareHouse/GetValue"
	WareHouse_DeleteValue_FullMethodName = "/warehouse.WareHouse/DeleteValue"
)

// WareHouseClient is the client API for WareHouse service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WareHouseClient interface {
	SetValue(ctx context.Context, in *Pair, opts ...grpc.CallOption) (*Empty, error)
	GetValue(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Result, error)
	DeleteValue(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Empty, error)
}

type wareHouseClient struct {
	cc grpc.ClientConnInterface
}

func NewWareHouseClient(cc grpc.ClientConnInterface) WareHouseClient {
	return &wareHouseClient{cc}
}

func (c *wareHouseClient) SetValue(ctx context.Context, in *Pair, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, WareHouse_SetValue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wareHouseClient) GetValue(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, WareHouse_GetValue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wareHouseClient) DeleteValue(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, WareHouse_DeleteValue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WareHouseServer is the server API for WareHouse service.
// All implementations must embed UnimplementedWareHouseServer
// for forward compatibility
type WareHouseServer interface {
	SetValue(context.Context, *Pair) (*Empty, error)
	GetValue(context.Context, *Key) (*Result, error)
	DeleteValue(context.Context, *Key) (*Empty, error)
	mustEmbedUnimplementedWareHouseServer()
}

// UnimplementedWareHouseServer must be embedded to have forward compatible implementations.
type UnimplementedWareHouseServer struct {
}

func (UnimplementedWareHouseServer) SetValue(context.Context, *Pair) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetValue not implemented")
}
func (UnimplementedWareHouseServer) GetValue(context.Context, *Key) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValue not implemented")
}
func (UnimplementedWareHouseServer) DeleteValue(context.Context, *Key) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteValue not implemented")
}
func (UnimplementedWareHouseServer) mustEmbedUnimplementedWareHouseServer() {}

// UnsafeWareHouseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WareHouseServer will
// result in compilation errors.
type UnsafeWareHouseServer interface {
	mustEmbedUnimplementedWareHouseServer()
}

func RegisterWareHouseServer(s grpc.ServiceRegistrar, srv WareHouseServer) {
	s.RegisterService(&WareHouse_ServiceDesc, srv)
}

func _WareHouse_SetValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pair)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WareHouseServer).SetValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WareHouse_SetValue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WareHouseServer).SetValue(ctx, req.(*Pair))
	}
	return interceptor(ctx, in, info, handler)
}

func _WareHouse_GetValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WareHouseServer).GetValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WareHouse_GetValue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WareHouseServer).GetValue(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _WareHouse_DeleteValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WareHouseServer).DeleteValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WareHouse_DeleteValue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WareHouseServer).DeleteValue(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

// WareHouse_ServiceDesc is the grpc.ServiceDesc for WareHouse service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WareHouse_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "warehouse.WareHouse",
	HandlerType: (*WareHouseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetValue",
			Handler:    _WareHouse_SetValue_Handler,
		},
		{
			MethodName: "GetValue",
			Handler:    _WareHouse_GetValue_Handler,
		},
		{
			MethodName: "DeleteValue",
			Handler:    _WareHouse_DeleteValue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "warehouse.proto",
}
