package repository

//go:generate mockgen -source=post_hot_repo.go -destination=../../internal/mocks/post_hot_repo_mock.go  -package mocks

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
	_tablePostHotName   = (&model.PostHotModel{}).TableName()
	_getPostHotSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetPostHotSQL = "SELECT * FROM %s WHERE id IN (%s)"
	_getHotPostListSQL  = "SELECT * FROM %s WHERE post_id <= ? order by score desc limit ?"
)

var _ PostHotRepo = (*postHotRepo)(nil)

// PostHotRepo define a repo interface
type PostHotRepo interface {
	CreatePostHot(ctx context.Context, db *gorm.DB, data *model.PostHotModel) (id int64, err error)
	UpdatePostHot(ctx context.Context, id int64, data *model.PostHotModel) error
	GetPostHot(ctx context.Context, id int64) (ret *model.PostHotModel, err error)
	BatchGetPostHot(ctx context.Context, ids []int64) (ret []*model.PostHotModel, err error)
	GetHotPostList(ctx context.Context, lastId int64, limit int32) (ret []*model.PostHotModel, err error)
}

type postHotRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
}

// NewPostHot new a repository and return
func NewPostHot(db *gorm.DB) PostHotRepo {
	return &postHotRepo{
		db:     db,
		tracer: otel.Tracer("postHotRepo"),
	}
}

// CreatePostHot create a item
func (r *postHotRepo) CreatePostHot(ctx context.Context, db *gorm.DB, data *model.PostHotModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create PostHot err")
	}

	return data.PostID, nil
}

// UpdatePostHot update item
func (r *postHotRepo) UpdatePostHot(ctx context.Context, id int64, data *model.PostHotModel) error {
	item, err := r.GetPostHot(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update PostHot err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// GetPostHot get a record
func (r *postHotRepo) GetPostHot(ctx context.Context, id int64) (ret *model.PostHotModel, err error) {
	data := new(model.PostHotModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getPostHotSQL, _tablePostHotName), id).Scan(&data).Error
	if err != nil {
		return
	}

	return data, nil
}

// BatchGetPostHot batch get items
func (r *postHotRepo) BatchGetPostHot(ctx context.Context, ids []int64) (ret []*model.PostHotModel, err error) {
	items := make([]*model.PostHotModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetPostHotSQL, _tablePostHotName), ids).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}

// GetHotPostList get hot post list
func (r *postHotRepo) GetHotPostList(ctx context.Context, lastId int64, limit int32) (ret []*model.PostHotModel, err error) {
	items := make([]*model.PostHotModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getHotPostListSQL, _tablePostHotName), lastId, limit).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}
