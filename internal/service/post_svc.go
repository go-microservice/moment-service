package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/utils"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	v1 "github.com/go-microservice/moment-service/api/moment/v1"
	"github.com/go-microservice/moment-service/internal/ecode"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"
	"github.com/go-microservice/moment-service/internal/tasks"
)

type PostType int
type DeleteType int
type VisibleType int

const (
	// PostTypeUnknown Post 类型
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
	content, err = getPostContent(postType, req)
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

	err = tasks.NewPublishPostTask(tasks.PublishPostPayload{
		PostID:    postID,
		AnchorUID: req.UserId,
	})
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
	pbPost, err := convertPost(*data)
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
	if len(req.Text) == 0 && len(req.Images) == 0 && req.Video != nil && len(req.Video.VideoKey) == 0 {
		return ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("param is empty"),
		})).Status(req).Err()
	}
	if len(req.Images) > 0 && req.Video != nil && len(req.Video.VideoKey) > 0 {
		return ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("only support for pic_keys or video_key"),
		})).Status(req).Err()
	}
	if req.Video != nil && len(req.Video.VideoKey) > 0 {
		if len(req.Video.VideoKey) == 0 || req.Video.Duration == 0 ||
			req.Video.Width == 0 || req.Video.Height == 0 {
			return ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": errors.New("video_duration or cover_key or width or height is empty"),
			})).Status(req).Err()
		}
	}

	return nil
}

func getPostType(req *v1.CreatePostRequest) PostType {
	if req.GetImages() != nil && len(req.GetImages()) > 0 {
		return PostTypeImage
	} else if req.GetVideo() != nil && len(req.GetVideo().VideoKey) > 0 {
		return PostTypeVideo
	} else if len(req.GetText()) > 0 {
		return PostTypeText
	}

	return PostTypeUnknown
}

func getPostContent(postType PostType, req *v1.CreatePostRequest) (string, error) {
	data := make(map[string]interface{})
	switch postType {
	case PostTypeText:
		data["text"] = req.GetText()
	case PostTypeImage:
		data["text"] = req.GetText()
		data["images"] = req.GetImages()
	case PostTypeVideo:
		video := req.Video
		data["text"] = req.GetText()
		data["video"] = map[string]interface{}{
			"video_key": video.GetVideoKey(),
			"duration":  video.GetDuration(),
			"cover_key": video.GetCoverKey(),
			"width":     video.GetWidth(),
			"height":    video.GetHeight(),
		}
	default:
		return "", ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("post_type is error"),
		})).Status(req).Err()
	}
	content, err := json.Marshal(data)
	if err != nil {
		return "", errors.Wrap(err, "getPostContent Marshal error")
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

	// check post if exist
	post, err := s.GetPost(ctx, &v1.GetPostRequest{Id: postID})
	if err != nil {
		return nil, err
	}

	// check if it has permission
	if req.GetUserId() != post.GetPost().UserId {
		return nil, ecode.ErrAccessDenied.WithDetails().Status(req).Err()
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
	pbPost, err := convertPost(*post)
	if err != nil {
		return nil, err
	}
	return &v1.GetPostReply{
		Post: pbPost,
	}, nil
}

func (s *PostServiceServer) BatchGetPosts(ctx context.Context, req *v1.BatchGetPostsRequest) (*v1.BatchGetPostsReply, error) {
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
		for {
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
		}
	}()

	for _, post := range posts {
		wg.Add(1)
		go func(info model.PostInfoModel) {
			defer func() {
				wg.Done()
				// catch error and avoid goroutine leak
				if r := recover(); r != nil {
					debugStr := utils.PrintStackTrace("BatchGetPosts: panic recovered", r)
					errChan <- fmt.Errorf("error: %s", debugStr)
				}
			}()

			mu.Lock()
			defer mu.Unlock()

			pbPost, err := convertPost(info)
			if err != nil {
				errChan <- err
			}

			m.Store(info.ID, pbPost)
		}(*post)

	}

	wg.Wait()
	close(errChan)
	close(finished)

	// 保证顺序
	for _, uid := range req.GetIds() {
		post, ok := m.Load(uid)
		if ok {
			pbPosts = append(pbPosts, post.(*v1.Post))
		}
	}

	return &v1.BatchGetPostsReply{
		Posts: pbPosts,
	}, nil
}

func (s *PostServiceServer) ListMyPosts(ctx context.Context, req *v1.ListMyPostsRequest) (*v1.ListMyPostsReply, error) {
	if req.GetPageToken() == 0 {
		req.PageToken = math.MaxInt64
	}
	if req.GetPageSize() == 0 {
		req.PageSize = 10
	}

	// get user posts
	userPosts, err := s.userPostRepo.GetUserPostByUserId(ctx, req.GetUserId(), cast.ToInt64(req.GetPageToken()), req.GetPageSize()+1)
	if err != nil {
		return nil, err
	}

	var (
		nextPageToken int64
	)
	if len(userPosts) > int(req.GetPageSize()) {
		nextPageToken = userPosts[len(userPosts)-1].ID
		userPosts = userPosts[:len(userPosts)-1]
	}

	// batch get post info
	var postIds []int64
	for _, userPost := range userPosts {
		postIds = append(postIds, userPost.PostID)
	}
	posts, err := s.BatchGetPosts(ctx, &v1.BatchGetPostsRequest{Ids: postIds})
	if err != nil {
		return nil, err
	}

	return &v1.ListMyPostsReply{
		Posts:         posts.GetPosts(),
		NextPageToken: nextPageToken,
	}, nil
}

func (s *PostServiceServer) ListLatestPosts(ctx context.Context, req *v1.ListLatestPostsRequest) (*v1.ListLatestPostsReply, error) {
	if req.GetPageToken() == 0 {
		req.PageToken = math.MaxInt64
	}
	if req.GetPageSize() == 0 {
		req.PageSize = 10
	}

	// get latest posts
	latestPosts, err := s.latestRepo.GetLatestPostList(ctx, cast.ToInt64(req.GetPageToken()), req.GetPageSize()+1)
	if err != nil {
		return nil, err
	}

	var (
		nextPageToken int64
	)
	if len(latestPosts) > int(req.GetPageSize()) {
		nextPageToken = latestPosts[len(latestPosts)-1].PostID
		latestPosts = latestPosts[:len(latestPosts)-1]
	}

	// batch get post info
	var postIds []int64
	for _, latestPost := range latestPosts {
		postIds = append(postIds, latestPost.PostID)
	}
	posts, err := s.BatchGetPosts(ctx, &v1.BatchGetPostsRequest{Ids: postIds})
	if err != nil {
		return nil, err
	}

	return &v1.ListLatestPostsReply{
		Posts:         posts.GetPosts(),
		NextPageToken: nextPageToken,
	}, nil
}

func (s *PostServiceServer) ListHotPosts(ctx context.Context, req *v1.ListHotPostsRequest) (*v1.ListHotPostsReply, error) {
	if req.GetPageToken() == 0 {
		req.PageToken = math.MaxInt64
	}
	if req.GetPageSize() == 0 {
		req.PageSize = 10
	}

	// get hot posts
	hotPosts, err := s.hotRepo.GetHotPostList(ctx, cast.ToInt64(req.GetPageToken()), req.GetPageSize()+1)
	if err != nil {
		return nil, err
	}

	var (
		nextPageToken int64
	)
	if len(hotPosts) > int(req.GetPageSize()) {
		nextPageToken = hotPosts[len(hotPosts)-1].PostID
		hotPosts = hotPosts[:len(hotPosts)-1]
	}

	// batch get post info
	var postIds []int64
	for _, hotPost := range hotPosts {
		postIds = append(postIds, hotPost.PostID)
	}
	posts, err := s.BatchGetPosts(ctx, &v1.BatchGetPostsRequest{Ids: postIds})
	if err != nil {
		return nil, err
	}

	return &v1.ListHotPostsReply{
		Posts:         posts.GetPosts(),
		NextPageToken: nextPageToken,
	}, nil
}

func convertPost(p model.PostInfoModel) (*v1.Post, error) {
	pbPost := &v1.Post{}
	err := copier.Copy(pbPost, &p)
	if err != nil {
		return nil, err
	}

	// NOTE: 字段大小写不一致时需要手动转换
	pbPost.Id = p.ID
	pbPost.Content, err = convertPostContent(p)
	if err != nil {
		return nil, err
	}
	return pbPost, nil
}

type PostText struct {
	Text string `json:"text"`
}

type PostImage struct {
	Text   string `json:"text"`
	Images []struct {
		ImageKey  string `json:"image_key"`
		ImageType string `json:"image_type"`
		Width     int32  `json:"width"`
		Height    int32  `json:"height"`
	} `json:"images"`
}

type PostVideo struct {
	Text        string `json:"text"`
	VideoKey    string `json:"video_key"`
	Duration    int    `json:"duration"`
	CoverKey    string `json:"cover_key"`
	CoverWidth  int    `json:"cover_width"`
	CoverHeight int    `json:"cover_height"`
}

// see: https://pkg.go.dev/google.golang.org/protobuf/types/known/structpb?utm_source=godoc#pkg-overview
func convertPostContent(p model.PostInfoModel) (ret *v1.Content, err error) {
	if len(p.Content) == 0 {
		return nil, nil
	}

	var content *v1.Content
	err = json.Unmarshal([]byte(p.Content), &content)
	if err != nil {
		return nil, errors.Wrap(err, "convertPostContent Unmarshal content error")
	}

	postType := PostType(p.PostType)
	switch postType {
	case PostTypeImage:
		for _, v := range content.GetImages() {
			v.ImageUrl = "http://aaa.cdnplus.com" + v.ImageKey
		}
	case PostTypeVideo:
		content.Video.CoverUrl = "http://aaa.cdnplus.com" + content.Video.CoverKey
	}

	return content, nil
}
