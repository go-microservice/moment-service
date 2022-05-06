package repository

//go:generate mockgen -source=comment_content_repo.go -destination=../../internal/mocks/comment_content_repo_mock.go  -package mocks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/go-microservice/moment-service/internal/cache"
	"github.com/go-microservice/moment-service/internal/model"
)

var (
	_tableCommentContentName   = (&model.CommentContentModel{}).TableName()
	_getCommentContentSQL      = "SELECT * FROM %s WHERE comment_id = ?"
	_batchGetCommentContentSQL = "SELECT * FROM %s WHERE comment_id IN (%s)"
)

var _ CommentContentRepo = (*commentContentRepo)(nil)

// CommentContentRepo define a repo interface
type CommentContentRepo interface {
	CreateCommentContent(ctx context.Context, db *gorm.DB, data *model.CommentContentModel) (id int64, err error)
	UpdateCommentContent(ctx context.Context, id int64, data *model.CommentContentModel) error
	GetCommentContent(ctx context.Context, id int64) (ret *model.CommentContentModel, err error)
	BatchGetCommentContent(ctx context.Context, ids []int64) (ret map[int64]*model.CommentContentModel, err error)
}

type commentContentRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.CommentContentCache
}

// NewCommentContent new a repository and return
func NewCommentContent(db *gorm.DB, cache cache.CommentContentCache) CommentContentRepo {
	return &commentContentRepo{
		db:     db,
		tracer: otel.Tracer("commentContent"),
		cache:  cache,
	}
}

// CreateCommentContent create a item
func (r *commentContentRepo) CreateCommentContent(ctx context.Context, db *gorm.DB, data *model.CommentContentModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create CommentContent err")
	}

	return data.CommentId, nil
}

// UpdateCommentContent update item
func (r *commentContentRepo) UpdateCommentContent(ctx context.Context, id int64, data *model.CommentContentModel) error {
	item, err := r.GetCommentContent(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update CommentContent err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelCommentContentCache(ctx, id)
	return nil
}

// GetCommentContent get a record
func (r *commentContentRepo) GetCommentContent(ctx context.Context, id int64) (ret *model.CommentContentModel, err error) {
	// read cache
	item, err := r.cache.GetCommentContentCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	// read db
	data := new(model.CommentContentModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getCommentContentSQL, _tableCommentContentName), id).Scan(&data).Error
	if err != nil {
		return
	}
	// write cache
	if data.CommentId > 0 {
		err = r.cache.SetCommentContentCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// BatchGetCommentContent batch get items
func (r *commentContentRepo) BatchGetCommentContent(ctx context.Context, ids []int64) (ret map[int64]*model.CommentContentModel, err error) {
	// read cache
	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGetCommentContentCache(ctx, ids)
	if err != nil {
		return nil, err
	}
	var (
		missedID []int64
		retMap   = make(map[int64]*model.CommentContentModel)
	)
	for _, v := range ids {
		item, ok := itemMap[v]
		if !ok {
			missedID = append(missedID, v)
			continue
		}
		retMap[v] = item
	}

	// get missed data
	if len(missedID) > 0 {
		var missedData []*model.CommentContentModel
		_sql := fmt.Sprintf(_batchGetCommentContentSQL, _tableCommentContentName, strings.Join(idsStr, ","))
		err = r.db.WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			for _, val := range missedData {
				retMap[val.CommentId] = val
			}
			err = r.cache.MultiSetCommentContentCache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return retMap, nil
}
