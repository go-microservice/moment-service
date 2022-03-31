package repository

//go:generate mockgen -source=user_post_repo.go -destination=../../internal/mocks/user_post_repo_mock.go  -package mocks

import (
	"context"
	"fmt"
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
)

var _ UserPostRepo = (*userPostRepo)(nil)

// UserPostRepo define a repo interface
type UserPostRepo interface {
	CreateUserPost(ctx context.Context, data *model.UserPostModel) (id int64, err error)
	UpdateUserPost(ctx context.Context, id int64, data *model.UserPostModel) error
	GetUserPost(ctx context.Context, id int64) (ret *model.UserPostModel, err error)
	BatchGetUserPost(ctx context.Context, ids []int64) (ret []*model.UserPostModel, err error)
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
func (r *userPostRepo) CreateUserPost(ctx context.Context, data *model.UserPostModel) (id int64, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
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
