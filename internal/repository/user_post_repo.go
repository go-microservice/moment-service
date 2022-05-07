package repository

//go:generate mockgen -source=user_post_repo.go -destination=../../internal/mocks/user_post_repo_mock.go  -package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/go-microservice/moment-service/internal/model"
)

var (
	_tableUserPostName   = (&model.UserPostModel{}).TableName()
	_getUserPostSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetUserPostSQL = "SELECT * FROM %s WHERE id IN (%s)"
	_getPostByUserIdSQL  = "SELECT * FROM %s WHERE user_id = ? and post_id <= ? order by post_id desc limit ?"
)

var _ UserPostRepo = (*userPostRepo)(nil)

// UserPostRepo define a repo interface
type UserPostRepo interface {
	CreateUserPost(ctx context.Context, db *gorm.DB, data *model.UserPostModel) (id int64, err error)
	UpdateUserPost(ctx context.Context, id int64, data *model.UserPostModel) error
	UpdateDelFlag(ctx context.Context, db *gorm.DB, id int64, delFlag int) error
	GetUserPost(ctx context.Context, id int64) (ret *model.UserPostModel, err error)
	BatchGetUserPost(ctx context.Context, ids []int64) (ret []*model.UserPostModel, err error)
	GetUserPostByUserId(ctx context.Context, userId int64, lastId int64, limit int32) (ret []*model.UserPostModel, err error)
}

type userPostRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
}

// NewUserPost new a repository and return
func NewUserPost(db *gorm.DB) UserPostRepo {
	return &userPostRepo{
		db:     db,
		tracer: otel.Tracer("userPostRepo"),
	}
}

// CreateUserPost create a item
func (r *userPostRepo) CreateUserPost(ctx context.Context, db *gorm.DB, data *model.UserPostModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create UserPost err")
	}

	return data.ID, nil
}

// UpdateUserPost update item
func (r *userPostRepo) UpdateUserPost(ctx context.Context, id int64, data *model.UserPostModel) error {
	item, err := r.GetUserPost(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update UserPost err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostRepo) UpdateDelFlag(ctx context.Context, db *gorm.DB, id int64, delFlag int) error {
	err := db.WithContext(ctx).Model(&model.UserPostModel{}).Where("post_id = ?", id).
		UpdateColumn("del_flag", delFlag).
		UpdateColumn("updated_at", time.Now().Unix()).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUserPost get a record
func (r *userPostRepo) GetUserPost(ctx context.Context, id int64) (ret *model.UserPostModel, err error) {
	data := new(model.UserPostModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getUserPostSQL, _tableUserPostName), id).Scan(&data).Error
	if err != nil {
		return
	}

	return data, nil
}

// BatchGetUserPost batch get items
func (r *userPostRepo) BatchGetUserPost(ctx context.Context, ids []int64) (ret []*model.UserPostModel, err error) {
	items := make([]*model.UserPostModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetUserPostSQL, _tableUserPostName), ids).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}

// GetUserPostByUserId get items by user id
func (r *userPostRepo) GetUserPostByUserId(ctx context.Context, userId int64, lastId int64, limit int32) (ret []*model.UserPostModel, err error) {
	items := make([]*model.UserPostModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getPostByUserIdSQL, _tableUserPostName), userId, lastId, limit).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}
