package repository

//go:generate mockgen -source=comment_info_repo.go -destination=../../internal/mocks/comment_info_repo_mock.go  -package mocks

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
	_tableCommentInfoName   = (&model.CommentInfoModel{}).TableName()
	_getCommentInfoSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetCommentInfoSQL = "SELECT * FROM %s WHERE id IN (%s)"
)

var _ CommentInfoRepo = (*commentInfoRepo)(nil)

// CommentInfoRepo define a repo interface
type CommentInfoRepo interface {
	CreateCommentInfo(ctx context.Context, db *gorm.DB, data *model.CommentInfoModel) (id int64, err error)
	UpdateCommentInfo(ctx context.Context, id int64, data *model.CommentInfoModel) error
	GetCommentInfo(ctx context.Context, id int64) (ret *model.CommentInfoModel, err error)
	BatchGetCommentInfo(ctx context.Context, ids []int64) (ret []*model.CommentInfoModel, err error)
}

type commentInfoRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.CommentInfoCache
}

// NewCommentInfo new a repository and return
func NewCommentInfo(db *gorm.DB, cache cache.CommentInfoCache) CommentInfoRepo {
	return &commentInfoRepo{
		db:     db,
		tracer: otel.Tracer("commentInfo"),
		cache:  cache,
	}
}

// CreateCommentInfo create a item
func (r *commentInfoRepo) CreateCommentInfo(ctx context.Context, db *gorm.DB, data *model.CommentInfoModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create CommentInfo err")
	}

	return data.ID, nil
}

// UpdateCommentInfo update item
func (r *commentInfoRepo) UpdateCommentInfo(ctx context.Context, id int64, data *model.CommentInfoModel) error {
	item, err := r.GetCommentInfo(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update CommentInfo err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelCommentInfoCache(ctx, id)
	return nil
}

// GetCommentInfo get a record
func (r *commentInfoRepo) GetCommentInfo(ctx context.Context, id int64) (ret *model.CommentInfoModel, err error) {
	// read cache
	item, err := r.cache.GetCommentInfoCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	// read db
	data := new(model.CommentInfoModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getCommentInfoSQL, _tableCommentInfoName), id).Scan(&data).Error
	if err != nil {
		return
	}
	// write cache
	if data.ID > 0 {
		err = r.cache.SetCommentInfoCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// BatchGetCommentInfo batch get items
func (r *commentInfoRepo) BatchGetCommentInfo(ctx context.Context, ids []int64) (ret []*model.CommentInfoModel, err error) {
	// read cache
	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGetCommentInfoCache(ctx, ids)
	if err != nil {
		return nil, err
	}
	var missedID []int64
	for _, v := range ids {
		item, ok := itemMap[cast.ToString(v)]
		if !ok {
			missedID = append(missedID, v)
			continue
		}
		ret = append(ret, item)
	}
	// get missed data
	if len(missedID) > 0 {
		var missedData []*model.CommentInfoModel
		_sql := fmt.Sprintf(_batchGetCommentInfoSQL, _tableCommentInfoName, strings.Join(idsStr, ","))
		err = r.db.WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSetCommentInfoCache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
}
