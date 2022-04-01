// Code generated protoc-gen-go-gin. DO NOT EDIT.
// protoc-gen-go-gin 0.0.5

package v1

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	app "github.com/go-eagle/eagle/pkg/app"
	errcode "github.com/go-eagle/eagle/pkg/errcode"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the eagle package it is being compiled against.

// context.
// metadata.
// gin.app.errcode.

var response = app.NewResponse()

type PostServiceHTTPServer interface {
	BatchGetPost(context.Context, *BatchGetPostRequest) (*BatchGetPostReply, error)
	CreatePost(context.Context, *CreatePostRequest) (*CreatePostReply, error)
	DeletePost(context.Context, *DeletePostRequest) (*DeletePostReply, error)
	GetPost(context.Context, *GetPostRequest) (*GetPostReply, error)
	ListHotPost(context.Context, *ListHotPostRequest) (*ListHotPostReply, error)
	ListLatestPost(context.Context, *ListLatestPostRequest) (*ListLatestPostReply, error)
	ListMyPost(context.Context, *ListMyPostRequest) (*ListMyPostReply, error)
	UpdatePost(context.Context, *UpdatePostRequest) (*UpdatePostReply, error)
}

func RegisterPostServiceHTTPServer(r gin.IRouter, srv PostServiceHTTPServer) {
	s := PostService{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

type PostService struct {
	server PostServiceHTTPServer
	router gin.IRouter
}

func (s *PostService) CreatePost_0(ctx *gin.Context) {
	var in CreatePostRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).CreatePost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) UpdatePost_0(ctx *gin.Context) {
	var in UpdatePostRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).UpdatePost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) DeletePost_0(ctx *gin.Context) {
	var in DeletePostRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).DeletePost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) GetPost_0(ctx *gin.Context) {
	var in GetPostRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).GetPost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) BatchGetPost_0(ctx *gin.Context) {
	var in BatchGetPostRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).BatchGetPost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) ListMyPost_0(ctx *gin.Context) {
	var in ListMyPostRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).ListMyPost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) ListLatestPost_0(ctx *gin.Context) {
	var in ListLatestPostRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).ListLatestPost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) ListHotPost_0(ctx *gin.Context) {
	var in ListHotPostRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		response.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(PostServiceHTTPServer).ListHotPost(newCtx, &in)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, out)
}

func (s *PostService) RegisterService() {
	s.router.Handle("POST", "post", s.CreatePost_0)
	s.router.Handle("PUT", "post", s.UpdatePost_0)
	s.router.Handle("DELETE", "post", s.DeletePost_0)
	s.router.Handle("GET", "post", s.GetPost_0)
	s.router.Handle("POST", "get/post", s.BatchGetPost_0)
	s.router.Handle("GET", "my/post", s.ListMyPost_0)
	s.router.Handle("GET", "latest/post", s.ListLatestPost_0)
	s.router.Handle("GET", "hot/post", s.ListHotPost_0)
}
