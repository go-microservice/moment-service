package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/go-eagle/eagle/pkg/log"

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

	postRepo       repository.PostInfoRepo
	cmtInfoRepo    repository.CommentInfoRepo
	cmtContentRepo repository.CommentContentRepo
	cmtLatestRepo  repository.CommentLatestRepo
	cmtHotRepo     repository.CommentHotRepo
}

func NewCommentServiceServer(
	postRepo repository.PostInfoRepo,
	cmtInfoRepo repository.CommentInfoRepo,
	cmtContentRepo repository.CommentContentRepo,
	cmtLatestRepo repository.CommentLatestRepo,
	cmtTopRepo repository.CommentHotRepo,
) *CommentServiceServer {
	return &CommentServiceServer{
		postRepo:       postRepo,
		cmtInfoRepo:    cmtInfoRepo,
		cmtContentRepo: cmtContentRepo,
		cmtLatestRepo:  cmtLatestRepo,
		cmtHotRepo:     cmtTopRepo,
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
		PostId:    req.GetPostId(),
		RootId:    req.GetRootId(),
		ParentId:  req.GetParentId(),
		UserId:    req.GetUserId(),
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
		CommentId:  cmtID,
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

	// create latest
	cmtLatest := &model.CommentLatestModel{
		CommentId: cmtID,
		PostID:    req.GetPostId(),
		RootID:    req.GetRootId(),
		ParentID:  req.GetParentId(),
		UserID:    req.GetUserId(),
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	_, err = s.cmtLatestRepo.CreateCommentLatest(ctx, tx, cmtLatest)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// create hot
	cmtHot := &model.CommentHotModel{
		CommentId: cmtID,
		PostID:    req.GetPostId(),
		RootID:    req.GetRootId(),
		ParentID:  req.GetParentId(),
		UserID:    req.GetUserId(),
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	_, err = s.cmtHotRepo.CreateCommentHot(ctx, tx, cmtHot)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// update comment count of post
	err = s.postRepo.IncrCommentCount(ctx, tx, req.GetPostId())
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

func (s *CommentServiceServer) ReplyComment(ctx context.Context, req *pb.ReplyCommentRequest) (*pb.ReplyCommentReply, error) {
	// check param
	if err := checkReplyParam(req); err != nil {
		return nil, err
	}

	// check comment if exist
	comment, err := s.GetComment(ctx, &pb.GetCommentRequest{Id: req.GetCommentId()})
	if err != nil {
		return nil, err
	}

	// start transaction
	tx := model.GetDB().Begin()
	if tx == nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	// create comment
	createTime := time.Now().Unix()
	cmtInfo := &model.CommentInfoModel{
		PostId:    comment.GetComment().GetPostId(),
		RootId:    comment.GetComment().GetId(),
		ParentId:  req.GetParentId(),
		UserId:    req.GetUserId(),
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
		CommentId:  cmtID,
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

	// update comment count of post
	err = s.postRepo.IncrCommentCount(ctx, tx, comment.GetComment().GetPostId())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// update reply count of comment
	err = s.cmtInfoRepo.IncrReplyCount(ctx, tx, comment.GetComment().GetId())
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

	return &pb.ReplyCommentReply{
		Comment: pbComment,
	}, nil
}

func (s *CommentServiceServer) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentReply, error) {
	if req.GetId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status().Err()
	}

	cmtInfo, err := s.cmtInfoRepo.GetCommentInfo(ctx, req.GetId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	infos, err := s.cmtInfoRepo.BatchGetCommentInfo(ctx, req.GetIds())
	if err != nil {
		return nil, ecode.ErrNotFound.WithDetails().Status().Err()
	}
	contents, err := s.cmtContentRepo.BatchGetCommentContent(ctx, req.GetIds())
	if err != nil {
		return nil, ecode.ErrNotFound.WithDetails().Status().Err()
	}

	var (
		comments []*pb.Comment
		m        sync.Map
		mu       sync.Mutex
	)

	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	go func() {
		select {
		case <-finished:
			return
		case err := <-errChan:
			if err != nil {
				// NOTE: if need, record log to file
			}
		case <-time.After(3 * time.Second):
			log.Warn(fmt.Errorf("list users timeout after 3 seconds"))
			return
		}
	}()

	for _, val := range infos {
		wg.Add(1)
		go func(info *model.CommentInfoModel) {
			defer func() {
				wg.Done()
			}()

			mu.Lock()
			defer mu.Unlock()

			content, ok := contents[info.ID]
			if !ok {
				return
			}
			// convert to pb comment
			pbComment, err := convertComment(info, content)
			if err != nil {
				return
			}

			m.Store(info.ID, pbComment)
		}(val)

	}

	wg.Wait()
	close(errChan)
	close(finished)

	// 保证顺序
	for _, uid := range req.GetIds() {
		comment, _ := m.Load(uid)
		comments = append(comments, comment.(*pb.Comment))
	}

	return &pb.BatchGetCommentReply{
		Comments: comments,
	}, nil
}

func (s *CommentServiceServer) ListHotComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentReply, error) {
	if req.GetPostId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}
	cmtIDs, err := s.cmtHotRepo.ListCommentHot(ctx, req.GetPostId(), req.GetLastId(), int(req.GetLimit())+1)
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

	items, err := s.BatchGetComment(ctx, &pb.BatchGetCommentRequest{Ids: cmtIDs})
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	return &pb.ListCommentReply{
		Items:   items.GetComments(),
		Count:   int64(len(items.GetComments())),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func (s *CommentServiceServer) ListLatestComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentReply, error) {
	if req.GetPostId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}
	cmtIDs, err := s.cmtLatestRepo.ListCommentLatest(ctx, req.GetPostId(), req.GetLastId(), int(req.GetLimit())+1)
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
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.PostId == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.UserId == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.Content == "" {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	return nil
}

func checkReplyParam(req *pb.ReplyCommentRequest) error {
	if req == nil {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.CommentId == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.UserId == 0 {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	if req.Content == "" {
		return ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
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
