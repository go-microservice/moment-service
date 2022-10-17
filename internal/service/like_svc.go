package service

import (
	"context"
	"math"
	"time"

	"github.com/jinzhu/copier"

	v1 "github.com/go-microservice/moment-service/api/moment/v1"
	"github.com/go-microservice/moment-service/internal/ecode"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"
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

	// LikeStatusDisliked 未点赞
	LikeStatusDisliked LikeStatus = 0
	// LikeStatusLiked 已点赞
	LikeStatusLiked LikeStatus = 1
)

var (
	_ v1.LikeServiceServer = (*LikeServiceServer)(nil)
)

type LikeServiceServer struct {
	v1.UnimplementedLikeServiceServer

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

func (s *LikeServiceServer) CreateLike(ctx context.Context, req *v1.CreateLikeRequest) (*v1.CreateLikeReply, error) {
	// check param
	if err := checkCreateLikeParam(req); err != nil {
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
	userLike, err := s.likeRepo.GetUserLike(ctx, req.GetUserId(), req.GetObjId(), req.GetObjType())
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}
	if hasLiked(userLike) {
		return &v1.CreateLikeReply{}, ecode.ErrSuccess.WithDetails().Status(req).Err()
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
	_, err = s.likeRepo.CreateUserLike(ctx, tx, likeData)
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

	return &v1.CreateLikeReply{}, nil
}
func (s *LikeServiceServer) UpdateLike(ctx context.Context, req *v1.UpdateLikeRequest) (*v1.UpdateLikeReply, error) {
	return &v1.UpdateLikeReply{}, nil
}
func (s *LikeServiceServer) DeleteLike(ctx context.Context, req *v1.DeleteLikeRequest) (*v1.DeleteLikeReply, error) {
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
		Status:    int(LikeStatusDisliked),
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
	return &v1.DeleteLikeReply{}, nil
}

func hasLiked(data *model.UserLikeModel) bool {
	if data == nil {
		return false
	}
	if data.ObjID > 0 && LikeStatus(data.Status) == LikeStatusLiked {
		return true
	}

	return false
}

func (s *LikeServiceServer) GetLike(ctx context.Context, req *v1.GetLikeRequest) (*v1.GetLikeReply, error) {
	return &v1.GetLikeReply{}, nil
}

func (s *LikeServiceServer) ListPostLike(ctx context.Context, req *v1.ListPostLikeRequest) (*v1.ListLikeReply, error) {
	if req.GetPostId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}
	likes, err := s.likeRepo.ListUserLikeByObj(ctx, int32(LikeTypePost), req.GetPostId(), req.GetLastId(), req.GetLimit()+1)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	var (
		hasMore bool
		lastId  int64
	)
	if len(likes) > int(req.GetLimit()) {
		hasMore = true
		lastId = likes[len(likes)-1].ID
		likes = likes[:len(likes)-1]
	}

	var items []*v1.Like
	for _, val := range likes {
		v, err := convertLike(val)
		if err != nil {
			continue
		}
		items = append(items, v)
	}

	return &v1.ListLikeReply{
		Items:   items,
		Count:   int64(len(likes)),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func (s *LikeServiceServer) ListCommentLike(ctx context.Context, req *v1.ListCommentLikeRequest) (*v1.ListLikeReply, error) {
	if req.GetCommentId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}
	likes, err := s.likeRepo.ListUserLikeByObj(ctx, int32(LikeTypeComment), req.GetCommentId(), req.GetLastId(), req.GetLimit()+1)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	var (
		hasMore bool
		lastId  int64
	)
	if len(likes) > int(req.GetLimit()) {
		hasMore = true
		lastId = likes[len(likes)-1].ID
		likes = likes[:len(likes)-1]
	}

	var items []*v1.Like
	for _, val := range likes {
		v, err := convertLike(val)
		if err != nil {
			continue
		}
		items = append(items, v)
	}

	return &v1.ListLikeReply{
		Items:   items,
		Count:   int64(len(likes)),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func convertLike(data *model.UserLikeModel) (*v1.Like, error) {
	pbLike := &v1.Like{}
	err := copier.Copy(pbLike, &data)
	if err != nil {
		return nil, err
	}

	// NOTE: 字段大小写不一致时需要手动转换
	pbLike.Id = data.ID
	pbLike.UserId = data.UserID

	return pbLike, nil
}

func checkCreateLikeParam(req *v1.CreateLikeRequest) error {
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

func checkDeleteLikeParam(req *v1.DeleteLikeRequest) error {
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
