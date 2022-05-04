package repository

//go:generate mockgen -source=comment_index_repo.go -destination=../../internal/mocks/comment_index_repo_mock.go  -package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/go-microservice/moment-service/internal/cache"
	"github.com/go-microservice/moment-service/internal/model"
)

type ListType int

const (
	CommentListTypeHot    ListType = 1
	CommentListTypeLatest ListType = 2
)

var (
	_tableCommentIndexName    = (&model.CommentIndexModel{}).TableName()
	_getLatestCommentIndexSQL = "SELECT id FROM %s WHERE obj_id = ? AND obj_type = ? and root_id=0 and id <= ? ORDER BY id DESC LIMIT ?"
	_getHotCommentIndexSQL    = "SELECT id FROM %s WHERE obj_id = ? AND obj_type = ? and root_id=0 and id <= ? ORDER BY score DESC LIMIT ?"
)

var _ CommentIndexRepo = (*commentIndexRepo)(nil)

// CommentIndexRepo define a repo interface
type CommentIndexRepo interface {
	CreateCommentIndex(ctx context.Context, db *gorm.DB, data *model.CommentIndexModel) (id int64, err error)
	UpdateDelFlag(ctx context.Context, db *gorm.DB, cmtID, objID int64, objType int, delFlag int) error
	UpdateScore(ctx context.Context, db *gorm.DB, cmtID, objID int64, objType int, score int) error
	GetLatestCommentIndex(ctx context.Context, objID int64, objType int, lastID int64, limit int) (ret []int64, err error)
	GetHotCommentIndex(ctx context.Context, objID int64, objType int, lastID int64, limit int) (ret []int64, err error)
}

type commentIndexRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.CommentIndexCache
}

// NewCommentIndex new a repository and return
func NewCommentIndex(db *gorm.DB, cache cache.CommentIndexCache) CommentIndexRepo {
	return &commentIndexRepo{
		db:     db,
		tracer: otel.Tracer("commentIndex"),
		cache:  cache,
	}
}

// CreateCommentIndex create a item
func (r *commentIndexRepo) CreateCommentIndex(ctx context.Context, db *gorm.DB, data *model.CommentIndexModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create CommentIndex err")
	}

	return data.Id, nil
}

// UpdateDelFlag update item
func (r *commentIndexRepo) UpdateDelFlag(ctx context.Context, db *gorm.DB, cmtID, objID int64, objType int, delFlag int) error {
	data := &model.CommentInfoModel{
		DelFlag:   delFlag,
		UpdatedAt: time.Now().Unix(),
	}
	err := db.Model(&model.CommentIndexModel{}).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelCommentIndexCache(ctx, objID, objType, int(CommentListTypeHot), cmtID)
	return nil
}

// UpdateScore update item
func (r *commentIndexRepo) UpdateScore(ctx context.Context, db *gorm.DB, cmtID, objID int64, objType int, score int) error {
	data := &model.CommentInfoModel{
		Score:     score,
		UpdatedAt: time.Now().Unix(),
	}
	err := db.Model(&model.CommentIndexModel{}).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelCommentIndexCache(ctx, objID, objType, int(CommentListTypeHot), cmtID)
	_ = r.cache.DelCommentIndexCache(ctx, objID, objType, int(CommentListTypeLatest), cmtID)
	return nil
}

func (r *commentIndexRepo) GetLatestCommentIndex(ctx context.Context, objID int64, objType int, lastID int64, limit int) (ret []int64, err error) {
	var (
		items []*model.CommentIndexModel
	)

	cmtIDs, err := r.cache.GetListCommentIndexCache(ctx, objID, objType, int(CommentListTypeLatest), lastID, limit)
	if err != nil {
		return nil, err
	}
	if len(cmtIDs) > 0 {
		return cmtIDs, nil
	}

	_sql := fmt.Sprintf(_getLatestCommentIndexSQL, _tableCommentIndexName)
	err = r.db.WithContext(ctx).Raw(_sql, objID, objType, lastID, limit).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return
	}

	for _, item := range items {
		ret = append(ret, item.Id)
	}

	return ret, nil
}

func (r *commentIndexRepo) GetHotCommentIndex(ctx context.Context, objID int64, objType int, lastID int64, limit int) (ret []int64, err error) {
	var (
		items []*model.CommentIndexModel
	)

	cmtIDs, err := r.cache.GetListCommentIndexCache(ctx, objID, objType, int(CommentListTypeHot), lastID, limit)
	if err != nil {
		return nil, err
	}
	if len(cmtIDs) > 0 {
		return cmtIDs, nil
	}

	_sql := fmt.Sprintf(_getHotCommentIndexSQL, _tableCommentIndexName)
	err = r.db.WithContext(ctx).Raw(_sql, objID, objType, lastID, limit).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return
	}

	for _, item := range items {
		ret = append(ret, item.Id)
	}

	return ret, nil
}
