package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-eagle/eagle/pkg/errcode"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	pb "github.com/go-microservice/moment-service/api/post/v1"
	"github.com/go-microservice/moment-service/internal/ecode"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"
)

type PostType int
type DeleteType int
type VisibleType int

const (
	PostTypeUnknown PostType = 0 // 未知
	PostTypeText    PostType = 1 // 文本
	PostTypeImage   PostType = 2 // 图片
	PostTypeVideo   PostType = 3 // 视频

	DelFlagNormal  DeleteType = 0 // 正常
	DelFlagByUser  DeleteType = 1 // 用户删除
	delFlagByAdmin DeleteType = 2 // 删除

	VisibleAll      VisibleType = 0 // 公开
	VisibleOnlySelf VisibleType = 1 // 仅自己可见
)

var (
	_ pb.PostServiceServer = (*PostServiceServer)(nil)
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewPostServiceServer)

type PostServiceServer struct {
	pb.UnimplementedPostServiceServer

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

func (s *PostServiceServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostReply, error) {
	// check param
	if err := checkParam(req); err != nil {
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
	_ = postID

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

	return &pb.CreatePostReply{
		Post: pbPost,
	}, nil
}

func checkParam(req *pb.CreatePostRequest) error {
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

func getPostType(req *pb.CreatePostRequest) PostType {
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

func getContent(postType PostType, req *pb.CreatePostRequest) (string, error) {
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

func (s *PostServiceServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostReply, error) {
	return &pb.UpdatePostReply{}, nil
}

func (s *PostServiceServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostReply, error) {
	return &pb.DeletePostReply{}, nil
}

func (s *PostServiceServer) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostReply, error) {
	if req.GetId() == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("post_id is empty"),
		})).Status(req).Err()
	}
	post, err := s.postRepo.GetPostInfo(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	pbPost, err := convertPost(post)
	if err != nil {
		return nil, err
	}
	return &pb.GetPostReply{
		Post: pbPost,
	}, nil
}

func (s *PostServiceServer) BatchGetPost(ctx context.Context, req *pb.BatchGetPostRequest) (*pb.BatchGetPostReply, error) {
	if len(req.GetIds()) == 0 {
		return nil, ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": errors.New("post_ids is empty"),
		})).Status(req).Err()
	}

	posts, err := s.postRepo.BatchGetPostInfo(ctx, req.GetIds())
	if err != nil {
		return nil, err
	}
	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, post := range posts {
		pbPost, err := convertPost(post)
		if err != nil {
			return nil, err
		}
		pbPosts = append(pbPosts, pbPost)
	}

	return &pb.BatchGetPostReply{
		Posts: pbPosts,
	}, nil
}

func (s *PostServiceServer) ListMyPost(ctx context.Context, req *pb.ListMyPostRequest) (*pb.ListMyPostReply, error) {
	return &pb.ListMyPostReply{}, nil
}

func (s *PostServiceServer) ListLatestPost(ctx context.Context, req *pb.ListLatestPostRequest) (*pb.ListLatestPostReply, error) {
	return &pb.ListLatestPostReply{}, nil
}

func (s *PostServiceServer) ListHotPost(ctx context.Context, req *pb.ListHotPostRequest) (*pb.ListHotPostReply, error) {
	return &pb.ListHotPostReply{}, nil
}

func convertPost(p *model.PostInfoModel) (*pb.Post, error) {
	pbPost := &pb.Post{}
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
	fmt.Println("------", rawContent)
	if err != nil {
		fmt.Println("--err1----", err)
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

	fmt.Println("data------", data)
	content, err := json.Marshal(data)
	if err != nil {
		fmt.Println("--err2----", err)
		return "", err
	}
	return string(content), nil
}
