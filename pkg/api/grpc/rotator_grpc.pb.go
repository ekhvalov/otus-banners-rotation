// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: v1/rotator.proto

package grpc

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

// RotatorClient is the client API for Rotator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RotatorClient interface {
	CreateBanner(ctx context.Context, in *CreateBannerRequest, opts ...grpc.CallOption) (*CreateBannerResponse, error)
	DeleteBanner(ctx context.Context, in *DeleteBannerRequest, opts ...grpc.CallOption) (*DeleteBannerResponse, error)
	CreateSlot(ctx context.Context, in *CreateSlotRequest, opts ...grpc.CallOption) (*CreateSlotResponse, error)
	DeleteSlot(ctx context.Context, in *DeleteSlotRequest, opts ...grpc.CallOption) (*DeleteSlotResponse, error)
	CreateSocialGroup(ctx context.Context, in *CreateSocialGroupRequest, opts ...grpc.CallOption) (*CreateSocialGroupResponse, error)
	DeleteSocialGroup(ctx context.Context, in *DeleteSocialGroupRequest, opts ...grpc.CallOption) (*DeleteSocialGroupResponse, error)
	AttachBanner(ctx context.Context, in *AttachBannerRequest, opts ...grpc.CallOption) (*AttachBannerResponse, error)
	DetachBanner(ctx context.Context, in *DetachBannerRequest, opts ...grpc.CallOption) (*DetachBannerResponse, error)
	ClickBanner(ctx context.Context, in *ClickBannerRequest, opts ...grpc.CallOption) (*ClickBannerResponse, error)
	SelectBanner(ctx context.Context, in *SelectBannerRequest, opts ...grpc.CallOption) (*SelectBannerResponse, error)
}

type rotatorClient struct {
	cc grpc.ClientConnInterface
}

func NewRotatorClient(cc grpc.ClientConnInterface) RotatorClient {
	return &rotatorClient{cc}
}

func (c *rotatorClient) CreateBanner(ctx context.Context, in *CreateBannerRequest, opts ...grpc.CallOption) (*CreateBannerResponse, error) {
	out := new(CreateBannerResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/CreateBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) DeleteBanner(ctx context.Context, in *DeleteBannerRequest, opts ...grpc.CallOption) (*DeleteBannerResponse, error) {
	out := new(DeleteBannerResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/DeleteBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) CreateSlot(ctx context.Context, in *CreateSlotRequest, opts ...grpc.CallOption) (*CreateSlotResponse, error) {
	out := new(CreateSlotResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/CreateSlot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) DeleteSlot(ctx context.Context, in *DeleteSlotRequest, opts ...grpc.CallOption) (*DeleteSlotResponse, error) {
	out := new(DeleteSlotResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/DeleteSlot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) CreateSocialGroup(ctx context.Context, in *CreateSocialGroupRequest, opts ...grpc.CallOption) (*CreateSocialGroupResponse, error) {
	out := new(CreateSocialGroupResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/CreateSocialGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) DeleteSocialGroup(ctx context.Context, in *DeleteSocialGroupRequest, opts ...grpc.CallOption) (*DeleteSocialGroupResponse, error) {
	out := new(DeleteSocialGroupResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/DeleteSocialGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) AttachBanner(ctx context.Context, in *AttachBannerRequest, opts ...grpc.CallOption) (*AttachBannerResponse, error) {
	out := new(AttachBannerResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/AttachBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) DetachBanner(ctx context.Context, in *DetachBannerRequest, opts ...grpc.CallOption) (*DetachBannerResponse, error) {
	out := new(DetachBannerResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/DetachBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) ClickBanner(ctx context.Context, in *ClickBannerRequest, opts ...grpc.CallOption) (*ClickBannerResponse, error) {
	out := new(ClickBannerResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/ClickBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rotatorClient) SelectBanner(ctx context.Context, in *SelectBannerRequest, opts ...grpc.CallOption) (*SelectBannerResponse, error) {
	out := new(SelectBannerResponse)
	err := c.cc.Invoke(ctx, "/otus.rotator.v1.Rotator/SelectBanner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RotatorServer is the server API for Rotator service.
// All implementations must embed UnimplementedRotatorServer
// for forward compatibility
type RotatorServer interface {
	CreateBanner(context.Context, *CreateBannerRequest) (*CreateBannerResponse, error)
	DeleteBanner(context.Context, *DeleteBannerRequest) (*DeleteBannerResponse, error)
	CreateSlot(context.Context, *CreateSlotRequest) (*CreateSlotResponse, error)
	DeleteSlot(context.Context, *DeleteSlotRequest) (*DeleteSlotResponse, error)
	CreateSocialGroup(context.Context, *CreateSocialGroupRequest) (*CreateSocialGroupResponse, error)
	DeleteSocialGroup(context.Context, *DeleteSocialGroupRequest) (*DeleteSocialGroupResponse, error)
	AttachBanner(context.Context, *AttachBannerRequest) (*AttachBannerResponse, error)
	DetachBanner(context.Context, *DetachBannerRequest) (*DetachBannerResponse, error)
	ClickBanner(context.Context, *ClickBannerRequest) (*ClickBannerResponse, error)
	SelectBanner(context.Context, *SelectBannerRequest) (*SelectBannerResponse, error)
	mustEmbedUnimplementedRotatorServer()
}

// UnimplementedRotatorServer must be embedded to have forward compatible implementations.
type UnimplementedRotatorServer struct {
}

func (UnimplementedRotatorServer) CreateBanner(context.Context, *CreateBannerRequest) (*CreateBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBanner not implemented")
}
func (UnimplementedRotatorServer) DeleteBanner(context.Context, *DeleteBannerRequest) (*DeleteBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBanner not implemented")
}
func (UnimplementedRotatorServer) CreateSlot(context.Context, *CreateSlotRequest) (*CreateSlotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSlot not implemented")
}
func (UnimplementedRotatorServer) DeleteSlot(context.Context, *DeleteSlotRequest) (*DeleteSlotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSlot not implemented")
}
func (UnimplementedRotatorServer) CreateSocialGroup(context.Context, *CreateSocialGroupRequest) (*CreateSocialGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSocialGroup not implemented")
}
func (UnimplementedRotatorServer) DeleteSocialGroup(context.Context, *DeleteSocialGroupRequest) (*DeleteSocialGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSocialGroup not implemented")
}
func (UnimplementedRotatorServer) AttachBanner(context.Context, *AttachBannerRequest) (*AttachBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttachBanner not implemented")
}
func (UnimplementedRotatorServer) DetachBanner(context.Context, *DetachBannerRequest) (*DetachBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetachBanner not implemented")
}
func (UnimplementedRotatorServer) ClickBanner(context.Context, *ClickBannerRequest) (*ClickBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClickBanner not implemented")
}
func (UnimplementedRotatorServer) SelectBanner(context.Context, *SelectBannerRequest) (*SelectBannerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SelectBanner not implemented")
}
func (UnimplementedRotatorServer) mustEmbedUnimplementedRotatorServer() {}

// UnsafeRotatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RotatorServer will
// result in compilation errors.
type UnsafeRotatorServer interface {
	mustEmbedUnimplementedRotatorServer()
}

func RegisterRotatorServer(s grpc.ServiceRegistrar, srv RotatorServer) {
	s.RegisterService(&Rotator_ServiceDesc, srv)
}

func _Rotator_CreateBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).CreateBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/CreateBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).CreateBanner(ctx, req.(*CreateBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_DeleteBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).DeleteBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/DeleteBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).DeleteBanner(ctx, req.(*DeleteBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_CreateSlot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSlotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).CreateSlot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/CreateSlot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).CreateSlot(ctx, req.(*CreateSlotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_DeleteSlot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSlotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).DeleteSlot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/DeleteSlot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).DeleteSlot(ctx, req.(*DeleteSlotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_CreateSocialGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSocialGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).CreateSocialGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/CreateSocialGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).CreateSocialGroup(ctx, req.(*CreateSocialGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_DeleteSocialGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSocialGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).DeleteSocialGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/DeleteSocialGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).DeleteSocialGroup(ctx, req.(*DeleteSocialGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_AttachBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttachBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).AttachBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/AttachBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).AttachBanner(ctx, req.(*AttachBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_DetachBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetachBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).DetachBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/DetachBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).DetachBanner(ctx, req.(*DetachBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_ClickBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClickBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).ClickBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/ClickBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).ClickBanner(ctx, req.(*ClickBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rotator_SelectBanner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SelectBannerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RotatorServer).SelectBanner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otus.rotator.v1.Rotator/SelectBanner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RotatorServer).SelectBanner(ctx, req.(*SelectBannerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Rotator_ServiceDesc is the grpc.ServiceDesc for Rotator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rotator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "otus.rotator.v1.Rotator",
	HandlerType: (*RotatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBanner",
			Handler:    _Rotator_CreateBanner_Handler,
		},
		{
			MethodName: "DeleteBanner",
			Handler:    _Rotator_DeleteBanner_Handler,
		},
		{
			MethodName: "CreateSlot",
			Handler:    _Rotator_CreateSlot_Handler,
		},
		{
			MethodName: "DeleteSlot",
			Handler:    _Rotator_DeleteSlot_Handler,
		},
		{
			MethodName: "CreateSocialGroup",
			Handler:    _Rotator_CreateSocialGroup_Handler,
		},
		{
			MethodName: "DeleteSocialGroup",
			Handler:    _Rotator_DeleteSocialGroup_Handler,
		},
		{
			MethodName: "AttachBanner",
			Handler:    _Rotator_AttachBanner_Handler,
		},
		{
			MethodName: "DetachBanner",
			Handler:    _Rotator_DetachBanner_Handler,
		},
		{
			MethodName: "ClickBanner",
			Handler:    _Rotator_ClickBanner_Handler,
		},
		{
			MethodName: "SelectBanner",
			Handler:    _Rotator_SelectBanner_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/rotator.proto",
}
