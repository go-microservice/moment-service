package service

import (
	"context"
	"time"

	"github.com/go-microservice/moment-service/internal/model"

	"github.com/go-microservice/moment-service/internal/ecode"

	"github.com/go-microservice/moment-service/internal/repository"

	pb "github.com/go-microservice/moment-service/api/like/v1"
)

type LikeType int
type LikeStatus int

const (
	// LikeTypeUnknown 未知类型
	LikeTypeUnknown LikeType = 0
	// LikeTypePost 帖子
	LikeTypePost LikeType = 1
	// LikeTypeComment 评论
	LikeTypeComment LikeType = 2

	// LikeStatusUnliked 未点赞
	LikeStatusUnliked LikeStatus = 0
	// LikeStatusLiked 已点赞
	LikeStatusLiked LikeStatus = 1
)

var (
	_ pb.LikeServiceServer = (*LikeServiceServer)(nil)
)

type LikeServiceServer struct {
	pb.UnimplementedLikeServiceServer

	likeRepo    repository.UserLikeRepo
	postRepo    repository.PostInfoRepo
	cmtInfoRepo repository.CommentInfoRepo
}

func NewLikeServiceServer(
	likeRepo repository.UserLikeRepo,
	postRepo repository.PostInfoRepo,
	cmtInfoRepo repository.CommentInfoRepo,
) *LikeServiceServer {
	return &LikeServiceServer{
		likeRepo:    likeRepo,
		postRepo:    postRepo,
		cmtInfoRepo: cmtInfoRepo,
	}
}

func (s *LikeServiceServer) CreateLike(ctx context.Context, req *pb.CreateLikeRequest) (*pb.CreateLikeReply, error) {
	// check param
	if err := checkCreateLikeParam(req); err != nil {
		return nil, err
	}

	// start transaction
	tx := model.GetDB().Begin()
	if tx == nil {
		return nil, ecode.ErrInternalError.WithDetails().Status().Err()
	}

	// create like
	likeData := &model.UserLikeModel{
		ObjType:   int64(req.GetObjType()),
		ObjID:     req.GetObjId(),
		UserID:    req.GetUserId(),
		Status:    int(LikeStatusLiked),
		CreatedAt: time.Now().Unix(),
	}
	_, err := s.likeRepo.CreateUserLike(ctx, tx, likeData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// update like count
	switch LikeType(req.GetObjType()) {
	case LikeTypePost:
		err = s.postRepo.IncrLikeCount(ctx, tx, req.GetObjId())
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	case LikeTypeComment:
		err = s.cmtInfoRepo.IncrLikeCount(ctx, tx, req.GetObjId())
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &pb.CreateLikeReply{}, nil
}
func (s *LikeServiceServer) UpdateLike(ctx context.Context, req *pb.UpdateLikeRequest) (*pb.UpdateLikeReply, error) {
	return &pb.UpdateLikeReply{}, nil
}
func (s *LikeServiceServer) DeleteLike(ctx context.Context, req *pb.DeleteLikeRequest) (*pb.DeleteLikeReply, error) {
	// check param
	if err := checkDeleteLikeParam(req); err != nil {
		return nil, err
	}

	// check object if exist
	switch LikeType(req.GetObjType()) {
	case LikeTypePost:
		post, err := s.postRepo.GetPostInfo(ctx, req.GetObjId())
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
		}
		if post == nil || post.ID == 0 {
			return nil, ecode.ErrNotFound.WithDetails().Status(req).Err()
		}
	case LikeTypeComment:
		cmt, err := s.cmtInfoRepo.GetCommentInfo(ctx, req.GetObjId())
		if err != nil {
			return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
		}
		if cmt == nil || cmt.ID == 0 {
			return nil, ecode.ErrNotFound.WithDetails().Status(req).Err()
		}
	}

	// check if liked

	// start transaction
	tx := model.GetDB().Begin()
	if tx == nil {
		return nil, ecode.ErrInternalError.WithDetails().Status().Err()
	}

	// create like
	likeData := &model.UserLikeModel{
		ObjType:   int64(req.GetObjType()),
		ObjID:     req.GetObjId(),
		UserID:    req.GetUserId(),
		Status:    int(LikeStatusUnliked),
		CreatedAt: time.Now().Unix(),
	}
	_, err := s.likeRepo.CreateUserLike(ctx, tx, likeData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// update like count
	switch LikeType(req.GetObjType()) {
	case LikeTypePost:
		err = s.postRepo.DecrLikeCount(ctx, tx, req.GetObjId())
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	case LikeTypeComment:
		err = s.cmtInfoRepo.DecrLikeCount(ctx, tx, req.GetObjId())
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &pb.DeleteLikeReply{}, nil
}
func (s *LikeServiceServer) GetLike(ctx context.Context, req *pb.GetLikeRequest) (*pb.GetLikeReply, error) {
	return &pb.GetLikeReply{}, nil
}
func (s *LikeServiceServer) ListLike(ctx context.Context, req *pb.ListLikeRequest) (*pb.ListLikeReply, error) {
	return &pb.ListLikeReply{}, nil
}

func checkCreateLikeParam(req *pb.CreateLikeRequest) error {
	if req == nil {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.GetObjType() == int32(LikeTypeUnknown) {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.GetObjId() == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.GetUserId() == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	return nil
}

func checkDeleteLikeParam(req *pb.DeleteLikeRequest) error {
	if req == nil {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.GetObjType() == int32(LikeTypeUnknown) {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.GetObjId() == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.GetUserId() == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	return nil
}
