package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	v1 "github.com/go-microservice/moment-service/api/moment/v1"

	"github.com/go-eagle/eagle/pkg/log"

	"github.com/go-microservice/moment-service/internal/ecode"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"
	"github.com/jinzhu/copier"
)

type CommentType int

const (
	// CommentTypeText text comment type
	CommentTypeText CommentType = 1
	// CommentTypeImage image comment type
	CommentTypeImage CommentType = 2
)

var (
	_ v1.CommentServiceServer = (*CommentServiceServer)(nil)
)

type CommentServiceServer struct {
	v1.UnimplementedCommentServiceServer

	postRepo       repository.PostInfoRepo
	cmtInfoRepo    repository.CommentInfoRepo
	cmtContentRepo repository.CommentContentRepo
	cmtLatestRepo  repository.CommentLatestRepo
	cmtHotRepo     repository.CommentHotRepo
	userCmtRepo    repository.UserCommentRepo
}

func NewCommentServiceServer(
	postRepo repository.PostInfoRepo,
	cmtInfoRepo repository.CommentInfoRepo,
	cmtContentRepo repository.CommentContentRepo,
	cmtLatestRepo repository.CommentLatestRepo,
	cmtTopRepo repository.CommentHotRepo,
	userCmtRepo repository.UserCommentRepo,
) *CommentServiceServer {
	return &CommentServiceServer{
		postRepo:       postRepo,
		cmtInfoRepo:    cmtInfoRepo,
		cmtContentRepo: cmtContentRepo,
		cmtLatestRepo:  cmtLatestRepo,
		cmtHotRepo:     cmtTopRepo,
		userCmtRepo:    userCmtRepo,
	}
}

func (s *CommentServiceServer) CreateComment(ctx context.Context, req *v1.CreateCommentRequest) (*v1.CreateCommentReply, error) {
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

	// user comment
	userComment := &model.UserCommentModel{
		CommentID: cmtID,
		UserID:    req.GetUserId(),
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	_, err = s.userCmtRepo.CreateUserComment(ctx, tx, userComment)
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

	return &v1.CreateCommentReply{
		Comment: pbComment,
	}, nil
}

func (s *CommentServiceServer) UpdateComment(ctx context.Context, req *v1.UpdateCommentRequest) (*v1.UpdateCommentReply, error) {
	return &v1.UpdateCommentReply{}, nil
}

func (s *CommentServiceServer) DeleteComment(ctx context.Context, req *v1.DeleteCommentRequest) (*v1.DeleteCommentReply, error) {
	if req.GetId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	commentID := req.GetId()

	// check comment if exist
	cmt, err := s.GetComment(ctx, &v1.GetCommentRequest{Id: commentID})
	if err != nil {
		return nil, err
	}

	// check if has delete permission
	if req.GetUserId() != cmt.GetComment().UserId {
		return nil, ecode.ErrAccessDenied.WithDetails().Status(req).Err()
	}

	// start transaction
	tx := model.GetDB().Begin()
	if tx == nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	err = s.cmtInfoRepo.UpdateDelFlag(ctx, tx, commentID, int(DelFlagByUser))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.cmtLatestRepo.UpdateDelFlag(ctx, tx, commentID, int(DelFlagByUser))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.cmtHotRepo.UpdateDelFlag(ctx, tx, commentID, int(DelFlagByUser))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.postRepo.DecrCommentCount(ctx, tx, cmt.GetComment().PostId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &v1.DeleteCommentReply{}, nil
}

func (s *CommentServiceServer) ReplyComment(ctx context.Context, req *v1.ReplyCommentRequest) (*v1.ReplyCommentReply, error) {
	// check param
	if err := checkReplyParam(req); err != nil {
		return nil, err
	}

	// check comment if exist
	comment, err := s.GetComment(ctx, &v1.GetCommentRequest{Id: req.GetCommentId()})
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

	// create latest for reply list
	cmtLatest := &model.CommentLatestModel{
		CommentId: cmtID,
		PostID:    comment.GetComment().GetPostId(),
		RootID:    req.GetCommentId(),
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

	return &v1.ReplyCommentReply{
		Comment: pbComment,
	}, nil
}

func (s *CommentServiceServer) GetComment(ctx context.Context, req *v1.GetCommentRequest) (*v1.GetCommentReply, error) {
	if req.GetId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	cmtInfo, err := s.cmtInfoRepo.GetCommentInfo(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if cmtInfo == nil || cmtInfo.ID == 0 {
		return nil, ecode.ErrNotFound.WithDetails().Status(req).Err()
	}
	cmtContent, err := s.cmtContentRepo.GetCommentContent(ctx, req.GetId())
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	// convert to pb comment
	pbComment, err := convertComment(cmtInfo, cmtContent)
	if err != nil {
		return nil, err
	}

	return &v1.GetCommentReply{
		Comment: pbComment,
	}, nil
}

func (s *CommentServiceServer) BatchGetComment(ctx context.Context, req *v1.BatchGetCommentsRequest) (*v1.BatchGetCommentsReply, error) {
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
		comments []*v1.Comment
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

	// keep order
	for _, uid := range req.GetIds() {
		comment, _ := m.Load(uid)
		comments = append(comments, comment.(*v1.Comment))
	}

	return &v1.BatchGetCommentsReply{
		Comments: comments,
	}, nil
}

func (s *CommentServiceServer) ListHotComment(ctx context.Context, req *v1.ListCommentsRequest) (*v1.ListCommentsReply, error) {
	if req.GetPostId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetPageToken() == 0 {
		req.PageToken = math.MaxInt64
	}
	if req.GetPageSize() == 0 {
		req.PageSize = 10
	}
	cmtIDs, err := s.cmtHotRepo.ListCommentHot(ctx, req.GetPostId(), req.GetPageToken(), int(req.GetPageSize())+1)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	var (
		nextPageToken int64
	)
	if len(cmtIDs) > int(req.GetPageSize()) {
		nextPageToken = cmtIDs[len(cmtIDs)-1]
		cmtIDs = cmtIDs[:len(cmtIDs)-1]
	}

	items, err := s.BatchGetComment(ctx, &v1.BatchGetCommentsRequest{Ids: cmtIDs})
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	return &v1.ListCommentsReply{
		Items:         items.GetComments(),
		Count:         int64(len(items.GetComments())),
		NextPageToken: nextPageToken,
	}, nil
}

func (s *CommentServiceServer) ListLatestComment(ctx context.Context, req *v1.ListCommentsRequest) (*v1.ListCommentsReply, error) {
	if req.GetPostId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetPageToken() == 0 {
		req.PageToken = math.MaxInt64
	}
	if req.GetPageSize() == 0 {
		req.PageSize = 10
	}
	cmtIDs, err := s.cmtLatestRepo.ListCommentLatest(ctx, req.GetPostId(), req.GetPageToken(), int(req.GetPageSize())+1)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	var (
		nextPageToken int64
	)
	if len(cmtIDs) > int(req.GetPageSize()) {
		nextPageToken = cmtIDs[len(cmtIDs)-1]
		cmtIDs = cmtIDs[:len(cmtIDs)-1]
	}

	cmts, err := s.BatchGetComment(ctx, &v1.BatchGetCommentsRequest{Ids: cmtIDs})
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	return &v1.ListCommentsReply{
		Items:         cmts.GetComments(),
		Count:         int64(len(cmts.GetComments())),
		NextPageToken: nextPageToken,
	}, nil
}

func (s *CommentServiceServer) ListReplyComment(ctx context.Context, req *v1.ListReplyCommentsRequest) (*v1.ListReplyCommentsReply, error) {
	if req.GetCommentId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}
	if req.GetPageToken() == 0 {
		req.PageToken = math.MaxInt64
	}
	if req.GetPageSize() == 0 {
		req.PageSize = 10
	}
	cmtIDs, err := s.cmtLatestRepo.ListReplyComment(ctx, req.GetCommentId(), req.GetPageToken(), int(req.GetPageSize())+1)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	var (
		nextPageToken int64
	)
	if len(cmtIDs) > int(req.GetPageSize()) {
		nextPageToken = cmtIDs[len(cmtIDs)-1]
		cmtIDs = cmtIDs[:len(cmtIDs)-1]
	}

	cmts, err := s.BatchGetComment(ctx, &v1.BatchGetCommentsRequest{Ids: cmtIDs})
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	return &v1.ListReplyCommentsReply{
		Items:         cmts.GetComments(),
		Count:         int64(len(cmts.GetComments())),
		NextPageToken: nextPageToken,
	}, nil
}

func checkCommentParam(req *v1.CreateCommentRequest) error {
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

func checkReplyParam(req *v1.ReplyCommentRequest) error {
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

func convertComment(cmt *model.CommentInfoModel, c *model.CommentContentModel) (*v1.Comment, error) {
	pbComment := &v1.Comment{}
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
