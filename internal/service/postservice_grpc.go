package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	v1 "github.com/go-microservice/moment-service/api/moment/v1"
	"github.com/go-microservice/moment-service/internal/ecode"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"
)

type PostType int
type DeleteType int
type VisibleType int

const (
	// Post 类型
	PostTypeUnknown PostType = 0 // 未知
	PostTypeText    PostType = 1 // 文本
	PostTypeImage   PostType = 2 // 图片
	PostTypeVideo   PostType = 3 // 视频

	DelFlagNormal  DeleteType = 0 // 正常
	DelFlagByUser  DeleteType = 1 // 用户删除
	DelFlagByAdmin DeleteType = 2 // 管理员删除

	VisibleAll      VisibleType = 0 // 公开
	VisibleOnlySelf VisibleType = 1 // 仅自己可见
)

var (
	_ v1.PostServiceServer = (*PostServiceServer)(nil)
)

type PostServiceServer struct {
	v1.UnimplementedPostServiceServer

	postRepo     repository.PostInfoRepo
	latestRepo   repository.PostLatestRepo
	hotRepo      repository.PostHotRepo
	userPostRepo repository.UserPostRepo
}

func NewPostServiceServer(
	postRepo repository.PostInfoRepo,
	latestRepo repository.PostLatestRepo,
	hotRepo repository.PostHotRepo,
	userPostRepo repository.UserPostRepo,
) *PostServiceServer {
	return &PostServiceServer{
		postRepo:     postRepo,
		latestRepo:   latestRepo,
		hotRepo:      hotRepo,
		userPostRepo: userPostRepo,
	}
}

func (s *PostServiceServer) CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.CreatePostReply, error) {
	// check param
	if err := checkPostParam(req); err != nil {
		return nil, err
	}

	var (
		err      error
		postType PostType
		content  string
	)
	postType = getPostType(req)
	content, err = getContent(postType, req)
	if err != nil {
		return nil, err
	}

	// start transaction
	tx := model.GetDB().Begin()
	if tx == nil {
		return nil, ecode.ErrInternalError.WithDetails().Status().Err()
	}

	// create post
	createTime := time.Now().Unix()
	data := &model.PostInfoModel{
		PostType:  int(postType),
		UserId:    req.UserId,
		Title:     req.Title,
		Content:   content,
		Longitude: float64(req.Longitude),
		Latitude:  float64(req.Latitude),
		Position:  req.Position,
		DelFlag:   int(DelFlagNormal),
		Visible:   int(VisibleOnlySelf),
		CreatedAt: createTime,
	}
	postID, err := s.postRepo.CreatePostInfo(ctx, tx, data)
	if err != nil {
		return nil, err
	}

	// create latest post
	latestData := &model.PostLatestModel{
		PostID:    postID,
		UserID:    req.UserId,
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	_, err = s.latestRepo.CreatePostLatest(ctx, tx, latestData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// create hot post
	hotData := &model.PostHotModel{
		PostID:    postID,
		UserID:    req.UserId,
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	_, err = s.hotRepo.CreatePostHot(ctx, tx, hotData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// create user post
	userPostData := &model.UserPostModel{
		UserID:    req.UserId,
		PostID:    postID,
		DelFlag:   int(DelFlagNormal),
		CreatedAt: createTime,
	}
	_, err = s.userPostRepo.CreateUserPost(ctx, tx, userPostData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	data.ID = postID
	// NOTE: 不能copy到嵌套的结构体中，所以单独出来copy
	pbPost, err := convertPost(data)
	if err != nil {
		return nil, err
	}

	return &v1.CreatePostReply{
		Post: pbPost,
	}, nil
}

func checkPostParam(req *v1.CreatePostRequest) error {
	if req.UserId == 0 {
		return ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("user_id is empty"),
		})).Status(req).Err()
	}
	if len(req.Text) == 0 && len(req.PicKeys) == 0 && len(req.VideoKey) == 0 {
		return ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("param is empty"),
		})).Status(req).Err()
	}
	if len(req.PicKeys) > 0 && len(req.VideoKey) > 0 {
		return ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("pic_keys and video_key is error"),
		})).Status(req).Err()
	}
	if len(req.VideoKey) > 0 {
		if len(req.CoverKey) == 0 || req.VideoDuration == 0 ||
			req.CoverWidth == 0 || req.CoverHeight == 0 {
			return ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": errors.New("video_duration or cover_key or width or height is empty"),
			})).Status(req).Err()
		}
	}

	return nil
}

func getPostType(req *v1.CreatePostRequest) PostType {
	if len(req.PicKeys) > 0 {
		return PostTypeImage
	}
	if len(req.VideoKey) > 0 {
		return PostTypeVideo
	}
	if len(req.Text) > 0 {
		return PostTypeText
	}

	return PostTypeUnknown
}

func getContent(postType PostType, req *v1.CreatePostRequest) (string, error) {
	data := make(map[string]interface{})
	switch postType {
	case PostTypeText:
		data["text"] = req.Text
	case PostTypeImage:
		// TODO: add width and height
		pics := strings.Split(req.PicKeys, ",")
		data["pic"] = pics
	case PostTypeVideo:
		data["video"] = map[string]interface{}{
			"video_key":    req.VideoKey,
			"duration":     req.VideoDuration,
			"cover_key":    req.CoverKey,
			"cover_width":  req.CoverWidth,
			"cover_height": req.CoverHeight,
		}
	default:
		return "", ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("post_type is error"),
		})).Status(req).Err()
	}
	content, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *PostServiceServer) UpdatePost(ctx context.Context, req *v1.UpdatePostRequest) (*v1.UpdatePostReply, error) {
	return &v1.UpdatePostReply{}, nil
}

func (s *PostServiceServer) DeletePost(ctx context.Context, req *v1.DeletePostRequest) (*v1.DeletePostReply, error) {
	if req.GetId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails().Status(req).Err()
	}

	postID := req.GetId()

	// check comment if exist
	_, err := s.GetPost(ctx, &v1.GetPostRequest{Id: postID})
	if err != nil {
		return nil, err
	}

	// start transaction
	tx := model.GetDB().Begin()
	if tx == nil {
		return nil, ecode.ErrInternalError.WithDetails().Status(req).Err()
	}

	err = s.postRepo.UpdateDelFlag(ctx, tx, postID, int(DelFlagByUser))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.latestRepo.UpdateDelFlag(ctx, tx, postID, int(DelFlagByUser))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.hotRepo.UpdateDelFlag(ctx, tx, postID, int(DelFlagByUser))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.userPostRepo.UpdateDelFlag(ctx, tx, postID, int(DelFlagByUser))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &v1.DeletePostReply{}, nil
}

func (s *PostServiceServer) GetPost(ctx context.Context, req *v1.GetPostRequest) (*v1.GetPostReply, error) {
	if req.GetId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("post_id is empty"),
		})).Status(req).Err()
	}
	post, err := s.postRepo.GetPostInfo(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	// check if post is exist
	if post == nil || post.ID == 0 {
		return nil, ecode.ErrNotFound.WithDetails().Status(req).Err()
	}
	pbPost, err := convertPost(post)
	if err != nil {
		return nil, err
	}
	return &v1.GetPostReply{
		Post: pbPost,
	}, nil
}

func (s *PostServiceServer) BatchGetPost(ctx context.Context, req *v1.BatchGetPostRequest) (*v1.BatchGetPostReply, error) {
	if len(req.GetIds()) == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("post_ids is empty"),
		})).Status(req).Err()
	}

	posts, err := s.postRepo.BatchGetPostInfo(ctx, req.GetIds())
	if err != nil {
		return nil, err
	}

	var (
		pbPosts []*v1.Post
		m       sync.Map
		mu      sync.Mutex
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

	for _, post := range posts {
		wg.Add(1)
		go func(info *model.PostInfoModel) {
			defer func() {
				wg.Done()
			}()

			mu.Lock()
			defer mu.Unlock()

			pbPost, err := convertPost(info)
			if err != nil {
				return
			}

			m.Store(info.ID, pbPost)
		}(post)

	}

	wg.Wait()
	close(errChan)
	close(finished)

	// 保证顺序
	for _, uid := range req.GetIds() {
		post, _ := m.Load(uid)
		pbPosts = append(pbPosts, post.(*v1.Post))
	}

	return &v1.BatchGetPostReply{
		Posts: pbPosts,
	}, nil
}

func (s *PostServiceServer) ListMyPost(ctx context.Context, req *v1.ListMyPostRequest) (*v1.ListMyPostReply, error) {
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}

	// get user posts
	userPosts, err := s.userPostRepo.GetUserPostByUserId(ctx, req.GetUserId(), req.GetLastId(), req.GetLimit()+1)
	if err != nil {
		return nil, err
	}

	var (
		hasMore bool
		lastId  int64
	)
	if len(userPosts) > int(req.GetLimit()) {
		hasMore = true
		lastId = userPosts[len(userPosts)-1].ID
		userPosts = userPosts[:len(userPosts)-1]
	}

	// batch get post info
	var postIds []int64
	for _, userPost := range userPosts {
		postIds = append(postIds, userPost.PostID)
	}
	posts, err := s.BatchGetPost(ctx, &v1.BatchGetPostRequest{Ids: postIds})
	if err != nil {
		return nil, err
	}

	return &v1.ListMyPostReply{
		Items:   posts.GetPosts(),
		Count:   int64(len(posts.GetPosts())),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func (s *PostServiceServer) ListLatestPost(ctx context.Context, req *v1.ListLatestPostRequest) (*v1.ListLatestPostReply, error) {
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}

	// get latest posts
	latestPosts, err := s.latestRepo.GetLatestPostList(ctx, req.GetLastId(), req.GetLimit()+1)
	if err != nil {
		return nil, err
	}

	var (
		hasMore bool
		lastId  int64
	)
	if len(latestPosts) > int(req.GetLimit()) {
		hasMore = true
		lastId = latestPosts[len(latestPosts)-1].PostID
		latestPosts = latestPosts[:len(latestPosts)-1]
	}

	// batch get post info
	var postIds []int64
	for _, latestPost := range latestPosts {
		postIds = append(postIds, latestPost.PostID)
	}
	posts, err := s.BatchGetPost(ctx, &v1.BatchGetPostRequest{Ids: postIds})
	if err != nil {
		return nil, err
	}

	return &v1.ListLatestPostReply{
		Items:   posts.GetPosts(),
		Count:   int64(len(posts.GetPosts())),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func (s *PostServiceServer) ListHotPost(ctx context.Context, req *v1.ListHotPostRequest) (*v1.ListHotPostReply, error) {
	if req.GetLastId() == 0 {
		req.LastId = math.MaxInt64
	}
	if req.GetLimit() == 0 {
		req.Limit = 10
	}

	// get hot posts
	hotPosts, err := s.hotRepo.GetHotPostList(ctx, req.GetLastId(), req.GetLimit()+1)
	if err != nil {
		return nil, err
	}

	var (
		hasMore bool
		lastId  int64
	)
	if len(hotPosts) > int(req.GetLimit()) {
		hasMore = true
		lastId = hotPosts[len(hotPosts)-1].PostID
		hotPosts = hotPosts[:len(hotPosts)-1]
	}

	// batch get post info
	var postIds []int64
	for _, hotPost := range hotPosts {
		postIds = append(postIds, hotPost.PostID)
	}
	posts, err := s.BatchGetPost(ctx, &v1.BatchGetPostRequest{Ids: postIds})
	if err != nil {
		return nil, err
	}

	return &v1.ListHotPostReply{
		Items:   posts.GetPosts(),
		Count:   int64(len(posts.GetPosts())),
		HasMore: hasMore,
		LastId:  lastId,
	}, nil
}

func convertPost(p *model.PostInfoModel) (*v1.Post, error) {
	pbPost := &v1.Post{}
	err := copier.Copy(pbPost, &p)
	if err != nil {
		return nil, err
	}

	// NOTE: 字段大小写不一致时需要手动转换
	pbPost.Id = p.ID
	pbPost.Content, err = convertContent(p)
	if err != nil {
		return nil, err
	}
	return pbPost, nil
}

func convertContent(p *model.PostInfoModel) (string, error) {
	if len(p.Content) == 0 {
		return "", nil
	}

	rawContent := make(map[string]interface{})
	err := json.Unmarshal([]byte(p.Content), &rawContent)
	if err != nil {
		return "", err
	}

	data := make(map[string]interface{})
	postType := PostType(p.PostType)
	switch postType {
	case PostTypeText:
		data["text"] = rawContent["text"]
	case PostTypeImage:
		data["pic"] = rawContent["pic"]
	case PostTypeVideo:
		vContent := rawContent["video"].(map[string]interface{})
		data["video"] = map[string]interface{}{
			"video_url":    vContent["video_key"],
			"duration":     vContent["duration"],
			"cover_url":    vContent["cover_key"],
			"cover_width":  vContent["cover_width"],
			"cover_height": vContent["cover_height"],
		}
	}

	content, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
