package repository

//go:generate mockgen -source=post_info_repo.go -destination=../../internal/mocks/post_info_repo_mock.go  -package mocks

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
	_tablePostInfoName   = (&model.PostInfoModel{}).TableName()
	_getPostInfoSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetPostInfoSQL = "SELECT * FROM %s WHERE id IN (%s)"
)

var _ PostInfoRepo = (*postInfoRepo)(nil)

// PostInfoRepo define a repo interface
type PostInfoRepo interface {
	CreatePostInfo(ctx context.Context, db *gorm.DB, data *model.PostInfoModel) (id int64, err error)
	UpdatePostInfo(ctx context.Context, id int64, data *model.PostInfoModel) error
	UpdateDelFlag(ctx context.Context, db *gorm.DB, id int64, delFlag int) error
	IncrCommentCount(ctx context.Context, db *gorm.DB, id int64) error
	DecrCommentCount(ctx context.Context, db *gorm.DB, id int64) error
	IncrLikeCount(ctx context.Context, db *gorm.DB, id int64) error
	DecrLikeCount(ctx context.Context, db *gorm.DB, id int64) error
	GetPostInfo(ctx context.Context, id int64) (ret *model.PostInfoModel, err error)
	BatchGetPostInfo(ctx context.Context, ids []int64) (ret []*model.PostInfoModel, err error)
}

type postInfoRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.PostInfoCache
}

// NewPostInfo new a repository and return
func NewPostInfo(db *gorm.DB, cache cache.PostInfoCache) PostInfoRepo {
	return &postInfoRepo{
		db:     db,
		tracer: otel.Tracer("postInfo"),
		cache:  cache,
	}
}

// CreatePostInfo create a item
func (r *postInfoRepo) CreatePostInfo(ctx context.Context, db *gorm.DB, data *model.PostInfoModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create PostInfo err")
	}

	return data.ID, nil
}

// UpdatePostInfo update item
func (r *postInfoRepo) UpdatePostInfo(ctx context.Context, id int64, data *model.PostInfoModel) error {
	item, err := r.GetPostInfo(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update PostInfo err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelPostInfoCache(ctx, id)
	return nil
}

func (r *postInfoRepo) UpdateDelFlag(ctx context.Context, db *gorm.DB, id int64, delFlag int) error {
	err := db.Model(&model.PostInfoModel{}).Where("id = ?", id).
		UpdateColumn("del_flag", delFlag).
		UpdateColumn("updated_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}
	// delete cache
	err = r.cache.DelPostInfoCache(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *postInfoRepo) IncrCommentCount(ctx context.Context, db *gorm.DB, id int64) error {
	err := db.Model(&model.PostInfoModel{}).Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
		UpdateColumn("updated_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}
	// delete cache
	err = r.cache.DelPostInfoCache(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *postInfoRepo) DecrCommentCount(ctx context.Context, db *gorm.DB, id int64) error {
	err := db.Model(&model.PostInfoModel{}).Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", -1)).
		UpdateColumn("updated_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}
	// delete cache
	err = r.cache.DelPostInfoCache(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *postInfoRepo) IncrLikeCount(ctx context.Context, db *gorm.DB, id int64) error {
	err := db.Model(&model.PostInfoModel{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).
		UpdateColumn("updated_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}
	// delete cache
	err = r.cache.DelPostInfoCache(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *postInfoRepo) DecrLikeCount(ctx context.Context, db *gorm.DB, id int64) error {
	// NOTE: run 2 sql, include: update field updated_at
	err := db.Model(&model.PostInfoModel{}).Where("id = ? AND like_count > 0", id).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).
		UpdateColumn("updated_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}
	// delete cache
	err = r.cache.DelPostInfoCache(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// GetPostInfo get a record
func (r *postInfoRepo) GetPostInfo(ctx context.Context, id int64) (ret *model.PostInfoModel, err error) {
	// read cache
	item, err := r.cache.GetPostInfoCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	// read db
	data := new(model.PostInfoModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getPostInfoSQL, _tablePostInfoName), id).Scan(&data).Error
	if err != nil {
		return
	}
	// write cache
	if data.ID > 0 {
		err = r.cache.SetPostInfoCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// BatchGetPostInfo batch get items
func (r *postInfoRepo) BatchGetPostInfo(ctx context.Context, ids []int64) (ret []*model.PostInfoModel, err error) {
	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGetPostInfoCache(ctx, ids)
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
		var missedData []*model.PostInfoModel
		_sql := fmt.Sprintf(_batchGetPostInfoSQL, _tablePostInfoName, strings.Join(idsStr, ","))
		err = r.db.WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSetPostInfoCache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
}
