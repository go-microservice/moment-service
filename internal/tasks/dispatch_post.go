package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/repository"

	"github.com/hibiken/asynq"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/pkg/errors"
)

const (
	// TypeDispatchPost task name
	TypeDispatchPost = "dispatch:post"
)

// DispatchPostPayload define data payload
type DispatchPostPayload struct {
	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id"`
}

// 主要用于将post分发给粉丝，最后批量写入db

// NewDispatchPostTask to create a task.
func NewDispatchPostTask(data DispatchPostPayload) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "[tasks] json marshal error, name: %s", TypeDispatchPost)
	}
	task := asynq.NewTask(TypeDispatchPost, payload)
	_, err = GetClient().Enqueue(task)
	if err != nil {
		return errors.Wrapf(err, "[tasks] Enqueue task error, name: %s", TypeDispatchPost)
	}

	return nil
}

// HandleDispatchPostTask to handle the input task.
func HandleDispatchPostTask(ctx context.Context, t *asynq.Task) error {
	var p DispatchPostPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		return err
	}

	// 将发布队列的数据写入到粉丝post
	repo := repository.NewUserPost(model.GetDB())
	data := &model.UserPostModel{
		UserID:    p.UserID,
		PostID:    p.PostID,
		DelFlag:   0,
		CreatedAt: time.Now().Unix(),
	}
	_, err := repo.CreateUserPost(ctx, model.GetDB(), data)
	if err != nil {
		log.Errorf("CreateUserPost failed, data: %v, error: %v", err)
		return err
	}

	return nil
}
