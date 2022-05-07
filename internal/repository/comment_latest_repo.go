package repository

//go:generate mockgen -source=comment_latest_repo.go -destination=../../internal/mocks/comment_latest_repo_mock.go  -package mocks

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

var (
	_tableCommentLatestName   = (&model.CommentLatestModel{}).TableName()
	_getCommentLatestSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetCommentLatestSQL = "SELECT * FROM %s WHERE id IN (%s)"
	_listCommentLatestSQL     = "SELECT comment_id FROM %s WHERE del_flag = 0 AND post_id = ? and root_id=0 and comment_id <= ? ORDER BY comment_id DESC LIMIT ?"
)

var _ CommentLatestRepo = (*commentLatestRepo)(nil)

// CommentLatestRepo define a repo interface
type CommentLatestRepo interface {
	CreateCommentLatest(ctx context.Context, db *gorm.DB, data *model.CommentLatestModel) (id int64, err error)
	UpdateCommentLatest(ctx context.Context, id int64, data *model.CommentLatestModel) error
	UpdateDelFlag(ctx context.Context, db *gorm.DB, id int64, delFlag int) error
	GetCommentLatest(ctx context.Context, id int64) (ret *model.CommentLatestModel, err error)
	ListCommentLatest(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error)
}

type commentLatestRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.CommentLatestCache
}

// NewCommentLatest new a repository and return
func NewCommentLatest(db *gorm.DB, cache cache.CommentLatestCache) CommentLatestRepo {
	return &commentLatestRepo{
		db:     db,
		tracer: otel.Tracer("commentLatest"),
		cache:  cache,
	}
}

// CreateCommentLatest create a item
func (r *commentLatestRepo) CreateCommentLatest(ctx context.Context, db *gorm.DB, data *model.CommentLatestModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create CommentLatest err")
	}

	return data.CommentId, nil
}

// UpdateCommentLatest update item
func (r *commentLatestRepo) UpdateCommentLatest(ctx context.Context, id int64, data *model.CommentLatestModel) error {
	item, err := r.GetCommentLatest(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update CommentLatest err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelCommentLatestCache(ctx, id)
	return nil
}

func (r *commentLatestRepo) UpdateDelFlag(ctx context.Context, db *gorm.DB, id int64, delFlag int) error {
	err := db.Model(&model.CommentLatestModel{}).Where("comment_id = ?", id).
		UpdateColumn("del_flag", delFlag).
		UpdateColumn("updated_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}
	// delete cache
	err = r.cache.DelCommentLatestCache(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// GetCommentLatest get a record
func (r *commentLatestRepo) GetCommentLatest(ctx context.Context, id int64) (ret *model.CommentLatestModel, err error) {
	// read cache
	item, err := r.cache.GetCommentLatestCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	// read db
	data := new(model.CommentLatestModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getCommentLatestSQL, _tableCommentLatestName), id).Scan(&data).Error
	if err != nil {
		return
	}
	// write cache
	if data.CommentId > 0 {
		err = r.cache.SetCommentLatestCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (r *commentLatestRepo) ListCommentLatest(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error) {
	var (
		items []*model.CommentLatestModel
	)

	cmtIDs, err := r.cache.GetListCommentLatestCache(ctx, postID, lastID, limit)
	if err != nil {
		return nil, err
	}
	if len(cmtIDs) > 0 {
		return cmtIDs, nil
	}

	_sql := fmt.Sprintf(_listCommentLatestSQL, _tableCommentLatestName)
	err = r.db.WithContext(ctx).Raw(_sql, postID, lastID, limit).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return
	}

	for _, item := range items {
		ret = append(ret, item.CommentId)
	}

	return ret, nil
}
