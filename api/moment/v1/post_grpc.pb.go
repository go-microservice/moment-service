// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.3
// source: api/moment/v1/post.proto

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

const (
	PostService_CreatePost_FullMethodName     = "/api.moment.v1.PostService/CreatePost"
	PostService_UpdatePost_FullMethodName     = "/api.moment.v1.PostService/UpdatePost"
	PostService_DeletePost_FullMethodName     = "/api.moment.v1.PostService/DeletePost"
	PostService_GetPost_FullMethodName        = "/api.moment.v1.PostService/GetPost"
	PostService_BatchGetPost_FullMethodName   = "/api.moment.v1.PostService/BatchGetPost"
	PostService_ListMyPost_FullMethodName     = "/api.moment.v1.PostService/ListMyPost"
	PostService_ListLatestPost_FullMethodName = "/api.moment.v1.PostService/ListLatestPost"
	PostService_ListHotPost_FullMethodName    = "/api.moment.v1.PostService/ListHotPost"
)

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	// 创建帖子
	CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostReply, error)
	// 更新帖子，暂时不提供
	UpdatePost(ctx context.Context, in *UpdatePostRequest, opts ...grpc.CallOption) (*UpdatePostReply, error)
	// 删除帖子
	DeletePost(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*DeletePostReply, error)
	// 根据id获取指定帖子
	GetPost(ctx context.Context, in *GetPostRequest, opts ...grpc.CallOption) (*GetPostReply, error)
	// 批量获取帖子
	BatchGetPost(ctx context.Context, in *BatchGetPostRequest, opts ...grpc.CallOption) (*BatchGetPostReply, error)
	// 我发布过的帖子列表
	ListMyPost(ctx context.Context, in *ListMyPostRequest, opts ...grpc.CallOption) (*ListMyPostReply, error)
	// 最新的帖子列表
	ListLatestPost(ctx context.Context, in *ListLatestPostRequest, opts ...grpc.CallOption) (*ListLatestPostReply, error)
	// 热门的帖子列表
	ListHotPost(ctx context.Context, in *ListHotPostRequest, opts ...grpc.CallOption) (*ListHotPostReply, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostReply, error) {
	out := new(CreatePostReply)
	err := c.cc.Invoke(ctx, PostService_CreatePost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) UpdatePost(ctx context.Context, in *UpdatePostRequest, opts ...grpc.CallOption) (*UpdatePostReply, error) {
	out := new(UpdatePostReply)
	err := c.cc.Invoke(ctx, PostService_UpdatePost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) DeletePost(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*DeletePostReply, error) {
	out := new(DeletePostReply)
	err := c.cc.Invoke(ctx, PostService_DeletePost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetPost(ctx context.Context, in *GetPostRequest, opts ...grpc.CallOption) (*GetPostReply, error) {
	out := new(GetPostReply)
	err := c.cc.Invoke(ctx, PostService_GetPost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) BatchGetPost(ctx context.Context, in *BatchGetPostRequest, opts ...grpc.CallOption) (*BatchGetPostReply, error) {
	out := new(BatchGetPostReply)
	err := c.cc.Invoke(ctx, PostService_BatchGetPost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) ListMyPost(ctx context.Context, in *ListMyPostRequest, opts ...grpc.CallOption) (*ListMyPostReply, error) {
	out := new(ListMyPostReply)
	err := c.cc.Invoke(ctx, PostService_ListMyPost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) ListLatestPost(ctx context.Context, in *ListLatestPostRequest, opts ...grpc.CallOption) (*ListLatestPostReply, error) {
	out := new(ListLatestPostReply)
	err := c.cc.Invoke(ctx, PostService_ListLatestPost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) ListHotPost(ctx context.Context, in *ListHotPostRequest, opts ...grpc.CallOption) (*ListHotPostReply, error) {
	out := new(ListHotPostReply)
	err := c.cc.Invoke(ctx, PostService_ListHotPost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations must embed UnimplementedPostServiceServer
// for forward compatibility
type PostServiceServer interface {
	// 创建帖子
	CreatePost(context.Context, *CreatePostRequest) (*CreatePostReply, error)
	// 更新帖子，暂时不提供
	UpdatePost(context.Context, *UpdatePostRequest) (*UpdatePostReply, error)
	// 删除帖子
	DeletePost(context.Context, *DeletePostRequest) (*DeletePostReply, error)
	// 根据id获取指定帖子
	GetPost(context.Context, *GetPostRequest) (*GetPostReply, error)
	// 批量获取帖子
	BatchGetPost(context.Context, *BatchGetPostRequest) (*BatchGetPostReply, error)
	// 我发布过的帖子列表
	ListMyPost(context.Context, *ListMyPostRequest) (*ListMyPostReply, error)
	// 最新的帖子列表
	ListLatestPost(context.Context, *ListLatestPostRequest) (*ListLatestPostReply, error)
	// 热门的帖子列表
	ListHotPost(context.Context, *ListHotPostRequest) (*ListHotPostReply, error)
	mustEmbedUnimplementedPostServiceServer()
}

// UnimplementedPostServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPostServiceServer struct {
}

func (UnimplementedPostServiceServer) CreatePost(context.Context, *CreatePostRequest) (*CreatePostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (UnimplementedPostServiceServer) UpdatePost(context.Context, *UpdatePostRequest) (*UpdatePostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePost not implemented")
}
func (UnimplementedPostServiceServer) DeletePost(context.Context, *DeletePostRequest) (*DeletePostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (UnimplementedPostServiceServer) GetPost(context.Context, *GetPostRequest) (*GetPostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPost not implemented")
}
func (UnimplementedPostServiceServer) BatchGetPost(context.Context, *BatchGetPostRequest) (*BatchGetPostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchGetPost not implemented")
}
func (UnimplementedPostServiceServer) ListMyPost(context.Context, *ListMyPostRequest) (*ListMyPostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMyPost not implemented")
}
func (UnimplementedPostServiceServer) ListLatestPost(context.Context, *ListLatestPostRequest) (*ListLatestPostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLatestPost not implemented")
}
func (UnimplementedPostServiceServer) ListHotPost(context.Context, *ListHotPostRequest) (*ListHotPostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListHotPost not implemented")
}
func (UnimplementedPostServiceServer) mustEmbedUnimplementedPostServiceServer() {}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_CreatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).CreatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_CreatePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).CreatePost(ctx, req.(*CreatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_UpdatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).UpdatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_UpdatePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).UpdatePost(ctx, req.(*UpdatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_DeletePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).DeletePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_DeletePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).DeletePost(ctx, req.(*DeletePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetPost(ctx, req.(*GetPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_BatchGetPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchGetPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).BatchGetPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_BatchGetPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).BatchGetPost(ctx, req.(*BatchGetPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_ListMyPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMyPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).ListMyPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_ListMyPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).ListMyPost(ctx, req.(*ListMyPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_ListLatestPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLatestPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).ListLatestPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_ListLatestPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).ListLatestPost(ctx, req.(*ListLatestPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_ListHotPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListHotPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).ListHotPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_ListHotPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).ListHotPost(ctx, req.(*ListHotPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.moment.v1.PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePost",
			Handler:    _PostService_CreatePost_Handler,
		},
		{
			MethodName: "UpdatePost",
			Handler:    _PostService_UpdatePost_Handler,
		},
		{
			MethodName: "DeletePost",
			Handler:    _PostService_DeletePost_Handler,
		},
		{
			MethodName: "GetPost",
			Handler:    _PostService_GetPost_Handler,
		},
		{
			MethodName: "BatchGetPost",
			Handler:    _PostService_BatchGetPost_Handler,
		},
		{
			MethodName: "ListMyPost",
			Handler:    _PostService_ListMyPost_Handler,
		},
		{
			MethodName: "ListLatestPost",
			Handler:    _PostService_ListLatestPost_Handler,
		},
		{
			MethodName: "ListHotPost",
			Handler:    _PostService_ListHotPost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/moment/v1/post.proto",
}
