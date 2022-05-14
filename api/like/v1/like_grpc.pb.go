// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.3
// source: api/like/v1/like.proto

package v1

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

// LikeServiceClient is the client API for LikeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LikeServiceClient interface {
	CreateLike(ctx context.Context, in *CreateLikeRequest, opts ...grpc.CallOption) (*CreateLikeReply, error)
	UpdateLike(ctx context.Context, in *UpdateLikeRequest, opts ...grpc.CallOption) (*UpdateLikeReply, error)
	DeleteLike(ctx context.Context, in *DeleteLikeRequest, opts ...grpc.CallOption) (*DeleteLikeReply, error)
	GetLike(ctx context.Context, in *GetLikeRequest, opts ...grpc.CallOption) (*GetLikeReply, error)
	ListPostLike(ctx context.Context, in *ListPostLikeRequest, opts ...grpc.CallOption) (*ListLikeReply, error)
	ListCommentLike(ctx context.Context, in *ListCommentLikeRequest, opts ...grpc.CallOption) (*ListLikeReply, error)
}

type likeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLikeServiceClient(cc grpc.ClientConnInterface) LikeServiceClient {
	return &likeServiceClient{cc}
}

func (c *likeServiceClient) CreateLike(ctx context.Context, in *CreateLikeRequest, opts ...grpc.CallOption) (*CreateLikeReply, error) {
	out := new(CreateLikeReply)
	err := c.cc.Invoke(ctx, "/api.like.v1.LikeService/CreateLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) UpdateLike(ctx context.Context, in *UpdateLikeRequest, opts ...grpc.CallOption) (*UpdateLikeReply, error) {
	out := new(UpdateLikeReply)
	err := c.cc.Invoke(ctx, "/api.like.v1.LikeService/UpdateLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) DeleteLike(ctx context.Context, in *DeleteLikeRequest, opts ...grpc.CallOption) (*DeleteLikeReply, error) {
	out := new(DeleteLikeReply)
	err := c.cc.Invoke(ctx, "/api.like.v1.LikeService/DeleteLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) GetLike(ctx context.Context, in *GetLikeRequest, opts ...grpc.CallOption) (*GetLikeReply, error) {
	out := new(GetLikeReply)
	err := c.cc.Invoke(ctx, "/api.like.v1.LikeService/GetLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) ListPostLike(ctx context.Context, in *ListPostLikeRequest, opts ...grpc.CallOption) (*ListLikeReply, error) {
	out := new(ListLikeReply)
	err := c.cc.Invoke(ctx, "/api.like.v1.LikeService/ListPostLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) ListCommentLike(ctx context.Context, in *ListCommentLikeRequest, opts ...grpc.CallOption) (*ListLikeReply, error) {
	out := new(ListLikeReply)
	err := c.cc.Invoke(ctx, "/api.like.v1.LikeService/ListCommentLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LikeServiceServer is the server API for LikeService service.
// All implementations must embed UnimplementedLikeServiceServer
// for forward compatibility
type LikeServiceServer interface {
	CreateLike(context.Context, *CreateLikeRequest) (*CreateLikeReply, error)
	UpdateLike(context.Context, *UpdateLikeRequest) (*UpdateLikeReply, error)
	DeleteLike(context.Context, *DeleteLikeRequest) (*DeleteLikeReply, error)
	GetLike(context.Context, *GetLikeRequest) (*GetLikeReply, error)
	ListPostLike(context.Context, *ListPostLikeRequest) (*ListLikeReply, error)
	ListCommentLike(context.Context, *ListCommentLikeRequest) (*ListLikeReply, error)
	mustEmbedUnimplementedLikeServiceServer()
}

// UnimplementedLikeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLikeServiceServer struct {
}

func (UnimplementedLikeServiceServer) CreateLike(context.Context, *CreateLikeRequest) (*CreateLikeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLike not implemented")
}
func (UnimplementedLikeServiceServer) UpdateLike(context.Context, *UpdateLikeRequest) (*UpdateLikeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLike not implemented")
}
func (UnimplementedLikeServiceServer) DeleteLike(context.Context, *DeleteLikeRequest) (*DeleteLikeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLike not implemented")
}
func (UnimplementedLikeServiceServer) GetLike(context.Context, *GetLikeRequest) (*GetLikeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLike not implemented")
}
func (UnimplementedLikeServiceServer) ListPostLike(context.Context, *ListPostLikeRequest) (*ListLikeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPostLike not implemented")
}
func (UnimplementedLikeServiceServer) ListCommentLike(context.Context, *ListCommentLikeRequest) (*ListLikeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCommentLike not implemented")
}
func (UnimplementedLikeServiceServer) mustEmbedUnimplementedLikeServiceServer() {}

// UnsafeLikeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LikeServiceServer will
// result in compilation errors.
type UnsafeLikeServiceServer interface {
	mustEmbedUnimplementedLikeServiceServer()
}

func RegisterLikeServiceServer(s grpc.ServiceRegistrar, srv LikeServiceServer) {
	s.RegisterService(&LikeService_ServiceDesc, srv)
}

func _LikeService_CreateLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).CreateLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.like.v1.LikeService/CreateLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).CreateLike(ctx, req.(*CreateLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_UpdateLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).UpdateLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.like.v1.LikeService/UpdateLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).UpdateLike(ctx, req.(*UpdateLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_DeleteLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).DeleteLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.like.v1.LikeService/DeleteLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).DeleteLike(ctx, req.(*DeleteLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_GetLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).GetLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.like.v1.LikeService/GetLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).GetLike(ctx, req.(*GetLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_ListPostLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPostLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).ListPostLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.like.v1.LikeService/ListPostLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).ListPostLike(ctx, req.(*ListPostLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_ListCommentLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCommentLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).ListCommentLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.like.v1.LikeService/ListCommentLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).ListCommentLike(ctx, req.(*ListCommentLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LikeService_ServiceDesc is the grpc.ServiceDesc for LikeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LikeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.like.v1.LikeService",
	HandlerType: (*LikeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLike",
			Handler:    _LikeService_CreateLike_Handler,
		},
		{
			MethodName: "UpdateLike",
			Handler:    _LikeService_UpdateLike_Handler,
		},
		{
			MethodName: "DeleteLike",
			Handler:    _LikeService_DeleteLike_Handler,
		},
		{
			MethodName: "GetLike",
			Handler:    _LikeService_GetLike_Handler,
		},
		{
			MethodName: "ListPostLike",
			Handler:    _LikeService_ListPostLike_Handler,
		},
		{
			MethodName: "ListCommentLike",
			Handler:    _LikeService_ListCommentLike_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/like/v1/like.proto",
}
