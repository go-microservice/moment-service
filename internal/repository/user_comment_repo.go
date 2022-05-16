package repository

//go:generate mockgen -source=user_comment_repo.go -destination=../../internal/mocks/user_comment_repo_mock.go  -package mocks

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
	_tableUserCommentName   = (&model.UserCommentModel{}).TableName()
	_getUserCommentSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetUserCommentSQL = "SELECT * FROM %s WHERE id IN (%s)"
)

var _ UserCommentRepo = (*userCommentRepo)(nil)

// UserCommentRepo define a repo interface
type UserCommentRepo interface {
	CreateUserComment(ctx context.Context, db *gorm.DB, data *model.UserCommentModel) (id int64, err error)
	UpdateUserComment(ctx context.Context, id int64, data *model.UserCommentModel) error
	GetUserComment(ctx context.Context, id int64) (ret *model.UserCommentModel, err error)
	BatchGetUserComment(ctx context.Context, ids []int64) (ret []*model.UserCommentModel, err error)
}

type userCommentRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.UserCommentCache
}

// NewUserComment new a repository and return
func NewUserComment(db *gorm.DB, cache cache.UserCommentCache) UserCommentRepo {
	return &userCommentRepo{
		db:     db,
		tracer: otel.Tracer("userComment"),
		cache:  cache,
	}
}

// CreateUserComment create a item
func (r *userCommentRepo) CreateUserComment(ctx context.Context, db *gorm.DB, data *model.UserCommentModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create UserComment err")
	}

	return data.ID, nil
}

// UpdateUserComment update item
func (r *userCommentRepo) UpdateUserComment(ctx context.Context, id int64, data *model.UserCommentModel) error {
	item, err := r.GetUserComment(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update UserComment err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelUserCommentCache(ctx, id)
	return nil
}

// GetUserComment get a record
func (r *userCommentRepo) GetUserComment(ctx context.Context, id int64) (ret *model.UserCommentModel, err error) {
	// read cache
	item, err := r.cache.GetUserCommentCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}
	// read db
	data := new(model.UserCommentModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserCommentSQL, _tableUserCommentName), id).Scan(&data).Error
	if err != nil {
		return
	}
	// write cache
	if data.ID > 0 {
		err = r.cache.SetUserCommentCache(ctx, id, data, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// BatchGetUserComment batch get items
func (r *userCommentRepo) BatchGetUserComment(ctx context.Context, ids []int64) (ret []*model.UserCommentModel, err error) {
	// read cache
	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGetUserCommentCache(ctx, ids)
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
		var missedData []*model.UserCommentModel
		_sql := fmt.Sprintf(_batchGetUserCommentSQL, _tableUserCommentName, strings.Join(idsStr, ","))
		err = r.db.WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSetUserCommentCache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
}
