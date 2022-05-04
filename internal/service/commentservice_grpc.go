package service

import (
	"context"
	"math"
	"time"

	pb "github.com/go-microservice/moment-service/api/comment/v1"
	"github.com/go-microservice/moment-service/internal/ecode"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"
	"github.com/jinzhu/copier"
)

type ObjType int

const (
	// post
	CommentObjTypePost ObjType = 1
)

var (
	_ pb.CommentServiceServer = (*CommentServiceServer)(nil)
)

type CommentServiceServer struct {
	pb.UnimplementedCommentServiceServer

	cmtInfoRepo    repository.CommentInfoRepo
	cmtContentRepo repository.CommentContentRepo
	cmtIndexRepo   repository.CommentIndexRepo
}

func NewCommentServiceServer(
	cmtInfoRepo repository.CommentInfoRepo,
	cmtContentRepo repository.CommentContentRepo,
	cmtIndexRepo repository.CommentIndexRepo,
) *CommentServiceServer {
	return &CommentServiceServer{
		cmtInfoRepo:    cmtInfoRepo,
		cmtContentRepo: cmtContentRepo,
		cmtIndexRepo:   cmtIndexRepo,
	}
}

func (s *CommentServiceServer) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentReply, error) {
	// check param
	if err := checkCommentParam(req); err != nil {
		return nil, err
	}

	// start transaction
	tx := model.GetDB().Begin()
	if tx == nil {
		return nil, ecode.ErrInternalError.WithDetails().Status().Err()
	}

	// create comment
	createTime := time.Now().Unix()
	cmtInfo := &model.CommentInfoModel{
		ObjType:   int(req.ObjType),
		ObjId:     req.ObjId,
		UserId:    req.UserId,
		RootId:    req.RootId,
		ParentId:  req.ParentId,
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	cmtID, err := s.cmtInfoRepo.CreateCommentInfo(ctx, tx, cmtInfo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// create content
	cmtContent := &model.CommentContentModel{
		Id:         cmtID,
		Content:    req.Content,
		DeviceType: req.DeviceType,
		IP:         req.Ip,
		CreatedAt:  createTime,
	}
	_, err = s.cmtContentRepo.CreateCommentContent(ctx, tx, cmtContent)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// create index
	cmtIndex := &model.CommentIndexModel{
		Id:        cmtID,
		ObjType:   int(req.ObjType),
		ObjId:     req.ObjId,
		RootId:    req.RootId,
		ParentId:  req.ParentId,
		UserId:    req.UserId,
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	_, err = s.cmtIndexRepo.CreateCommentIndex(ctx, tx, cmtIndex)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	cmtInfo.ID = cmtID

	// convert to pb comment
	pbComment, err := convertComment(cmtInfo, cmtContent)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCommentReply{
		Comment: pbComment,
	}, nil
}

func (s *CommentServiceServer) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentReply, error) {
	return &pb.UpdateCommentReply{}, nil
}

func (s *CommentServiceServer) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentReply, error) {
	return &pb.DeleteCommentReply{}, nil
}

func (s *CommentServiceServer) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentReply, error) {
	if req.GetId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status().Err()
	}

	cmtInfo, err := s.cmtInfoRepo.GetCommentInfo(ctx, req.GetId())
	if err != nil {
		return nil, ecode.ErrNotFound.WithDetails().Status().Err()
	}
	cmtContent, err := s.cmtContentRepo.GetCommentContent(ctx, req.GetId())
	if err != nil {
		return nil, ecode.ErrNotFound.WithDetails().Status().Err()
	}

	// convert to pb comment
	pbComment, err := convertComment(cmtInfo, cmtContent)
	if err != nil {
		return nil, err
	}

	return &pb.GetCommentReply{
		Comment: pbComment,
	}, nil
}

func (s *CommentServiceServer) BatchGetComment(ctx context.Context, req *pb.BatchGetCommentRequest) (*pb.BatchGetCommentReply, error) {
	if len(req.GetIds()) == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	cmtInfos, err := s.cmtInfoRepo.BatchGetCommentInfo(ctx, req.GetIds())
	if err != nil {
		return nil, ecode.ErrNotFound.WithDetails().Status().Err()
	}
	cmtContents, err := s.cmtContentRepo.BatchGetCommentContent(ctx, req.GetIds())
	if err != nil {
		return nil, ecode.ErrNotFound.WithDetails().Status().Err()
	}

	var comments []*pb.Comment
	for key, _ := range cmtInfos {
		// convert to pb comment
		pbComment, err := convertComment(cmtInfos[key], cmtContents[key])
		if err != nil {
			continue
		}
		comments = append(comments, pbComment)
	}

	return &pb.BatchGetCommentReply{
		Comments: comments,
	}, nil
}

func (s *CommentServiceServer) ListHotComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentReply, error) {
	if req.GetObjId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}
	cmtIDs, err := s.cmtIndexRepo.GetHotCommentIndex(ctx, req.GetObjId(), int(CommentObjTypePost), req.GetLastId(), int(req.GetLimit()))
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	var (
		hasMore bool
		lastId  int64
	)
	if len(cmtIDs) > int(req.GetLimit()) {
		hasMore = true
		lastId = cmtIDs[len(cmtIDs)-1]
		cmtIDs = cmtIDs[:len(cmtIDs)-1]
	}

	cmts, err := s.BatchGetComment(ctx, &pb.BatchGetCommentRequest{Ids: cmtIDs})
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	return &pb.ListCommentReply{
		Items:   cmts.GetComments(),
		Count:   int64(len(cmts.GetComments())),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func (s *CommentServiceServer) ListLatestComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentReply, error) {
	if req.GetObjId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}
	cmtIDs, err := s.cmtIndexRepo.GetLatestCommentIndex(ctx, req.GetObjId(), int(CommentObjTypePost), req.GetLastId(), int(req.GetLimit()))
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	var (
		hasMore bool
		lastId  int64
	)
	if len(cmtIDs) > int(req.GetLimit()) {
		hasMore = true
		lastId = cmtIDs[len(cmtIDs)-1]
		cmtIDs = cmtIDs[:len(cmtIDs)-1]
	}

	cmts, err := s.BatchGetComment(ctx, &pb.BatchGetCommentRequest{Ids: cmtIDs})
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	return &pb.ListCommentReply{
		Items:   cmts.GetComments(),
		Count:   int64(len(cmts.GetComments())),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func checkCommentParam(req *pb.CreateCommentRequest) error {
	if req == nil {
		return ecode.ErrInvalidArgument.WithDetails().Status().Err()
	}

	if req.ObjType == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status().Err()
	}

	if req.ObjId == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status().Err()
	}

	if req.UserId == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status().Err()
	}

	if req.Content == "" {
		return ecode.ErrInvalidArgument.WithDetails().Status().Err()
	}

	return nil
}

func convertComment(cmt *model.CommentInfoModel, c *model.CommentContentModel) (*pb.Comment, error) {
	pbComment := &pb.Comment{}
	err := copier.Copy(pbComment, &cmt)
	if err != nil {
		return nil, err
	}

	// NOTE: 字段大小写不一致时需要手动转换
	pbComment.Id = cmt.ID
	pbComment.Content = c.Content
	pbComment.DeviceType = c.DeviceType
	pbComment.Ip = c.IP

	return pbComment, nil
}
