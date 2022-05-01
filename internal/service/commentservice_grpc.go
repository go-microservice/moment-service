package service

import (
	"context"
	"time"

	pb "github.com/go-microservice/moment-service/api/comment/v1"
	"github.com/go-microservice/moment-service/internal/ecode"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"
	"github.com/jinzhu/copier"
)

var (
	_ pb.CommentServiceServer = (*CommentServiceServer)(nil)
)

type CommentServiceServer struct {
	pb.UnimplementedCommentServiceServer

	cmtInfoRepo    repository.CommentInfoRepo
	cmtContentRepo repository.CommentContentRepo
}

func NewCommentServiceServer(
	cmtInfoRepo repository.CommentInfoRepo,
	cmtContentRepo repository.CommentContentRepo,
) *CommentServiceServer {
	return &CommentServiceServer{
		cmtInfoRepo:    cmtInfoRepo,
		cmtContentRepo: cmtContentRepo,
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

	cmtContent := &model.CommentContentModel{
		CommentID:  cmtID,
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
	return &pb.GetCommentReply{}, nil
}

func (s *CommentServiceServer) BatchGetComment(ctx context.Context, req *pb.BatchGetCommentRequest) (*pb.BatchGetCommentReply, error) {
	return &pb.BatchGetCommentReply{}, nil
}

func (s *CommentServiceServer) ListComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentReply, error) {
	return &pb.ListCommentReply{}, nil
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
