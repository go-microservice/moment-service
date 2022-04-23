package repository

//go:generate mockgen -source=post_latest_repo.go -destination=../../internal/mocks/post_latest_repo_mock.go  -package mocks

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
	_tablePostLatestName   = (&model.PostLatestModel{}).TableName()
	_getPostLatestSQL      = "SELECT * FROM %s WHERE id = ?"
	_batchGetPostLatestSQL = "SELECT * FROM %s WHERE id IN (%s)"
	_getLatestPostListSQL  = "SELECT * FROM %s WHERE post_id <= ? order by post_id desc limit ?"
)

var _ PostLatestRepo = (*postLatestRepo)(nil)

// PostLatestRepo define a repo interface
type PostLatestRepo interface {
	CreatePostLatest(ctx context.Context, db *gorm.DB, data *model.PostLatestModel) (id int64, err error)
	UpdatePostLatest(ctx context.Context, id int64, data *model.PostLatestModel) error
	GetPostLatest(ctx context.Context, id int64) (ret *model.PostLatestModel, err error)
	BatchGetPostLatest(ctx context.Context, ids []int64) (ret []*model.PostLatestModel, err error)
	GetLatestPostList(ctx context.Context, lastId int64, limit int32) (ret []*model.PostLatestModel, err error)
}

type postLatestRepo struct {
	db     *gorm.DB
	tracer trace.Tracer
}

// NewPostLatest new a repository and return
func NewPostLatest(db *gorm.DB) PostLatestRepo {
	return &postLatestRepo{
		db:     db,
		tracer: otel.Tracer("postLatestRepo"),
	}
}

// CreatePostLatest create a item
func (r *postLatestRepo) CreatePostLatest(ctx context.Context, db *gorm.DB, data *model.PostLatestModel) (id int64, err error) {
	err = db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo] create PostLatest err")
	}

	return data.PostID, nil
}

// UpdatePostLatest update item
func (r *postLatestRepo) UpdatePostLatest(ctx context.Context, id int64, data *model.PostLatestModel) error {
	item, err := r.GetPostLatest(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[repo] update PostLatest err: %v", err)
	}
	err = r.db.Model(&item).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// GetPostLatest get a record
func (r *postLatestRepo) GetPostLatest(ctx context.Context, id int64) (ret *model.PostLatestModel, err error) {
	data := new(model.PostLatestModel)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getPostLatestSQL, _tablePostLatestName), id).Scan(&data).Error
	if err != nil {
		return
	}

	return data, nil
}

// BatchGetPostLatest batch get items
func (r *postLatestRepo) BatchGetPostLatest(ctx context.Context, ids []int64) (ret []*model.PostLatestModel, err error) {
	items := make([]*model.PostLatestModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_batchGetPostLatestSQL, _tablePostLatestName), ids).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}

// GetLatestPostList get latest post list
func (r *postLatestRepo) GetLatestPostList(ctx context.Context, lastId int64, limit int32) (ret []*model.PostLatestModel, err error) {
	items := make([]*model.PostLatestModel, 0)
	err = r.db.WithContext(ctx).Raw(fmt.Sprintf(_getLatestPostListSQL, _tablePostLatestName), lastId, limit).Scan(&items).Error
	if err != nil {
		return
	}
	return items, nil
}
