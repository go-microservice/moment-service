package tasks

import (
	"context"
	"encoding/json"

	"github.com/go-microservice/moment-service/internal/repository"
	relationV1 "github.com/go-microservice/relation-service/api/relation/v1"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

const (
	// TypePublishPost task name
	TypePublishPost = "publish:post"
)

// PublishPostPayload define data payload
type PublishPostPayload struct {
	PostID    int64 `json:"post_id"`
	AnchorUID int64 `json:"anchor_uid"`
}

// 主要用于将新发布的post写入发布队列，后用task将拉取粉丝列表该post分发给粉丝队列

// NewPublishPostTask to create a task.
func NewPublishPostTask(data PublishPostPayload) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "[tasks] json marshal error, name: %s", TypePublishPost)
	}
	task := asynq.NewTask(TypePublishPost, payload)
	_, err = GetClient().Enqueue(task)
	if err != nil {
		return errors.Wrapf(err, "[tasks] Enqueue task error, name: %s", TypePublishPost)
	}

	return nil
}

// HandlePublishPostTask to handle the input task.
func HandlePublishPostTask(ctx context.Context, t *asynq.Task) error {
	var p PublishPostPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		return err
	}

	relationClient := repository.NewRelationClient()

	var (
		hasMore int32
		lastId  int64
		limit   int32 = 20
	)

	// 从关系服务获取粉丝用户，然后写入到分发队列
	for {
		// 1. 拉取用户的粉丝列表
		followerUIDs := make([]int64, 0)
		in := &relationV1.FollowerListRequest{
			UserId: p.AnchorUID,
			LastId: lastId,
			Limit:  limit + 1,
		}
		ret, err := relationClient.GetFollowerList(ctx, in)
		if err != nil {
			log.Errorf("GetFollowerList in: %v, error: %v: %w", in, err, asynq.SkipRetry)
			return err
		}

		relationRet := ret.GetResult()
		if int32(len(relationRet)) > limit {
			hasMore = 1
			lastId = relationRet[len(relationRet)-1].Id
			relationRet = relationRet[0 : len(relationRet)-1]
		}
		for _, v := range relationRet {
			followerUIDs = append(followerUIDs, v.GetFollowerUid())
		}

		// 2. 遍历粉丝数据并写入到分发队列里
		for _, uid := range followerUIDs {
			data := DispatchPostPayload{
				UserID: uid,
				PostID: p.PostID,
			}
			err := NewDispatchPostTask(data)
			if err != nil {
				log.Errorf("NewDispatchPostTask failed, err:%v, data: %w", err, asynq.SkipRetry)
				continue
			}
		}

		if hasMore == 0 {
			break
		}
	}

	return nil
}
