package repository

//go:generate mockgen -source=comment_hot_repo.go -destination=../../internal/mocks/comment_hot_repo_mock.go  -package mocks

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
	_tableCommentHotName = (&model.CommentHotModel{}).TableName()
	_getCommentHotSQL    = "SELECT * FROM %s WHERE comment_id = ?"
	_listCommentHotSQL   = "SELECT comment_id FROM %s WHERE del_flag = 0 AND post_id = ? and root_id=0 and comment_id <= ? ORDER BY score DESC LIMIT ?"
)

var _ CommentHotRepo = (*commentHotRepo)(nil)

// CommentHotRepo define a repo interface
type CommentHotRepo interface {
	CreateCommentHot(ctx context.Context, db *gorm.DB, data *model.CommentHotModel) (id int64, err error)
	UpdateCommentHot(ctx context.Context, id int64, data *model.CommentHotModel) error
	GetCommentHot(ctx context.Context, id int64) (ret *model.CommentHotModel, err error)
	ListCommentHot(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error)
}

type commentHotRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.CommentHotCache
}

// NewCommentHot new a repository and return
func NewCommentHot(db *gorm.DB, cache cache.CommentHotCache) CommentHotRepo {
	return &commentHotRepo{
		db:     db,
		tracer: otel.Tracer("commentHot"),
		cache:  cache,
	}
}

// CreateCommentHot create a item
func (r *commentHotRepo) CreateCommentHot(ctx context.Context, db *gorm.DB, data *model.CommentHotModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create CommentHot err")
	}

	return data.CommentId, nil
}

// UpdateCommentHot update item
func (r *commentHotRepo) UpdateCommentHot(ctx context.Context, id int64, data *model.CommentHotModel) error {
	item, err := r.GetCommentHot(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update CommentHot err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelCommentHotCache(ctx, id)
	return nil
}

// GetCommentHot get a record
func (r *commentHotRepo) GetCommentHot(ctx context.Context, id int64) (ret *model.CommentHotModel, err error) {
	// read cache
	item, err := r.cache.GetCommentHotCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	// read db
	data := new(model.CommentHotModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getCommentHotSQL, _tableCommentHotName), id).Scan(&data).Error
	if err != nil {
		return
	}
	// write cache
	if data.CommentId > 0 {
		err = r.cache.SetCommentHotCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (r *commentHotRepo) ListCommentHot(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error) {
	var (
		items []*model.CommentHotModel
	)

	cmtIDs, err := r.cache.GetListCommentHotCache(ctx, postID, lastID, limit)
	if err != nil {
		return nil, err
	}
	if len(cmtIDs) > 0 {
		return cmtIDs, nil
	}

	_sql := fmt.Sprintf(_listCommentHotSQL, _tableCommentHotName)
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
