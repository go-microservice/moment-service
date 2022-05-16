package repository

//go:generate mockgen -source=user_like_repo.go -destination=../../internal/mocks/user_like_repo_mock.go  -package mocks

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
	_tableUserLikeName    = (&model.UserLikeModel{}).TableName()
	_createSQL            = "INSERT IGNORE INTO %s SET obj_type=?, obj_id=?, user_id=?, status=1, created_at=? ON duplicate key update status=?"
	_getUserLikeSQL       = "SELECT user_id, obj_type, obj_id, status FROM %s WHERE user_id=? AND obj_type=? AND obj_id=?"
	_batchGetUserLikeSQL  = "SELECT * FROM %s WHERE id IN (%s)"
	_listUserLikeByObjSQL = "SELECT * FROM %s WHERE obj_type=? AND obj_id=? AND status=1 and id <=? ORDER BY id DESC limit ?"
)

var _ UserLikeRepo = (*userLikeRepo)(nil)

// UserLikeRepo define a repo interface
type UserLikeRepo interface {
	CreateUserLike(ctx context.Context, db *gorm.DB, data *model.UserLikeModel) (id int64, err error)
	UpdateUserLike(ctx context.Context, id int64, data *model.UserLikeModel) error
	GetUserLike(ctx context.Context, userID, objID int64, objType int32) (ret *model.UserLikeModel, err error)
	BatchGetUserLike(ctx context.Context, ids []int64) (ret []*model.UserLikeModel, err error)
	ListUserLikeByObj(ctx context.Context, objType int32, objID, lastID int64, limit int32) (ret []*model.UserLikeModel, err error)
}

type userLikeRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
	cache  cache.UserLikeCache
}

// NewUserLike new a repository and return
func NewUserLike(db *gorm.DB, cache cache.UserLikeCache) UserLikeRepo {
	return &userLikeRepo{
		db:     db,
		tracer: otel.Tracer("userLike"),
		cache:  cache,
	}
}

// CreateUserLike create a item
func (r *userLikeRepo) CreateUserLike(ctx context.Context, db *gorm.DB, data *model.UserLikeModel) (id int64, err error) {
	sql := fmt.Sprintf(_createSQL, data.TableName())
	err = db.WithContext(ctx).Exec(sql, data.ObjType, data.ObjID, data.UserID, data.CreatedAt, data.Status).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create UserLike err")
	}

	return data.ID, nil
}

// UpdateUserLike update item
func (r *userLikeRepo) UpdateUserLike(ctx context.Context, id int64, data *model.UserLikeModel) error {
	item := &model.UserLikeModel{}
	err := r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	// delete cache
	_ = r.cache.DelUserLikeCache(ctx, id)
	return nil
}

// GetUserLike get a record
// TODO: add cache
func (r *userLikeRepo) GetUserLike(ctx context.Context, userID, objID int64, objType int32) (ret *model.UserLikeModel, err error) {
	// read db
	data := new(model.UserLikeModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserLikeSQL, _tableUserLikeName), userID, objType, objID).Scan(&data).Error
	if err != nil {
		return
	}
	return data, nil
}

// BatchGetUserLike batch get items
func (r *userLikeRepo) BatchGetUserLike(ctx context.Context, ids []int64) (ret []*model.UserLikeModel, err error) {
	// read cache
	idsStr := cast.ToStringSlice(ids)
	itemMap, err := r.cache.MultiGetUserLikeCache(ctx, ids)
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
		var missedData []*model.UserLikeModel
		_sql := fmt.Sprintf(_batchGetUserLikeSQL, _tableUserLikeName, strings.Join(idsStr, ","))
		err = r.db.WithContext(ctx).Raw(_sql).Scan(&missedData).Error
		if err != nil {
			// you can degrade to ignore error
			return nil, err
		}
		if len(missedData) > 0 {
			ret = append(ret, missedData...)
			err = r.cache.MultiSetUserLikeCache(ctx, missedData, 5*time.Minute)
			if err != nil {
				// you can degrade to ignore error
				return nil, err
			}
		}
	}
	return ret, nil
}

func (r *userLikeRepo) ListUserLikeByObj(ctx context.Context, objType int32, objID, lastID int64, limit int32) (ret []*model.UserLikeModel, err error) {
	var items []*model.UserLikeModel
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_listUserLikeByObjSQL, _tableUserLikeName), objType, objID, lastID, limit).Scan(&items).Error
	if err != nil {
		return
	}

	return items, nil
}
