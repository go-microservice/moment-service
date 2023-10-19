package cache

//go:generate mockgen -source=internal/cache/user_comment_cache.go -destination=internal/mock/user_comment_cache_mock.go  -package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	redis "github.com/redis/go-redis/v9"

	"github.com/go-microservice/moment-service/internal/model"
)

const (
	// PrefixUserCommentCacheKey cache prefix
	PrefixUserCommentCacheKey = "user:comment:%d"
)

// UserCommentCache define cache interface
type UserCommentCache interface {
	SetUserCommentCache(ctx context.Context, id int64, data *model.UserCommentModel, duration time.Duration) error
	GetUserCommentCache(ctx context.Context, id int64) (data *model.UserCommentModel, err error)
	MultiGetUserCommentCache(ctx context.Context, ids []int64) (map[string]*model.UserCommentModel, error)
	MultiSetUserCommentCache(ctx context.Context, data []*model.UserCommentModel, duration time.Duration) error
	DelUserCommentCache(ctx context.Context, id int64) error
}

// userCommentCache define cache struct
type userCommentCache struct {
	cache cache.Cache
}

// NewUserCommentCache new a cache
func NewUserCommentCache(rdb *redis.Client) UserCommentCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &userCommentCache{
		cache: cache.NewRedisCache(rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserCommentModel{}
		}),
	}
}

// GetUserCommentCacheKey get cache key
func (c *userCommentCache) GetUserCommentCacheKey(id int64) string {
	return fmt.Sprintf(PrefixUserCommentCacheKey, id)
}

// SetUserCommentCache write to cache
func (c *userCommentCache) SetUserCommentCache(ctx context.Context, id int64, data *model.UserCommentModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserCommentCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserCommentCache get from cache
func (c *userCommentCache) GetUserCommentCache(ctx context.Context, id int64) (data *model.UserCommentModel, err error) {
	cacheKey := c.GetUserCommentCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserCommentCache batch get cache
func (c *userCommentCache) MultiGetUserCommentCache(ctx context.Context, ids []int64) (map[string]*model.UserCommentModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserCommentCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.UserCommentModel)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSetUserCommentCache batch set cache
func (c *userCommentCache) MultiSetUserCommentCache(ctx context.Context, data []*model.UserCommentModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserCommentCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelUserCommentCache delete cache
func (c *userCommentCache) DelUserCommentCache(ctx context.Context, id int64) error {
	cacheKey := c.GetUserCommentCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *userCommentCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetUserCommentCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
