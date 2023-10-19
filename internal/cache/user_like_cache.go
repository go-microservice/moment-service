package cache

//go:generate mockgen -source=internal/cache/user_like_cache.go -destination=internal/mock/user_like_cache_mock.go  -package mock

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
	// PrefixUserLikeCacheKey cache prefix
	PrefixUserLikeCacheKey = "user:like:%d"
)

// UserLikeCache define cache interface
type UserLikeCache interface {
	SetUserLikeCache(ctx context.Context, id int64, data *model.UserLikeModel, duration time.Duration) error
	GetUserLikeCache(ctx context.Context, id int64) (data *model.UserLikeModel, err error)
	MultiGetUserLikeCache(ctx context.Context, ids []int64) (map[string]*model.UserLikeModel, error)
	MultiSetUserLikeCache(ctx context.Context, data []*model.UserLikeModel, duration time.Duration) error
	DelUserLikeCache(ctx context.Context, id int64) error
}

// userLikeCache define cache struct
type userLikeCache struct {
	cache cache.Cache
}

// NewUserLikeCache new a cache
func NewUserLikeCache(rdb *redis.Client) UserLikeCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &userLikeCache{
		cache: cache.NewRedisCache(rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserLikeModel{}
		}),
	}
}

// GetUserLikeCacheKey get cache key
func (c *userLikeCache) GetUserLikeCacheKey(id int64) string {
	return fmt.Sprintf(PrefixUserLikeCacheKey, id)
}

// SetUserLikeCache write to cache
func (c *userLikeCache) SetUserLikeCache(ctx context.Context, id int64, data *model.UserLikeModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserLikeCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserLikeCache get from cache
func (c *userLikeCache) GetUserLikeCache(ctx context.Context, id int64) (data *model.UserLikeModel, err error) {
	cacheKey := c.GetUserLikeCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserLikeCache batch get cache
func (c *userLikeCache) MultiGetUserLikeCache(ctx context.Context, ids []int64) (map[string]*model.UserLikeModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserLikeCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.UserLikeModel)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSetUserLikeCache batch set cache
func (c *userLikeCache) MultiSetUserLikeCache(ctx context.Context, data []*model.UserLikeModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserLikeCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelUserLikeCache delete cache
func (c *userLikeCache) DelUserLikeCache(ctx context.Context, id int64) error {
	cacheKey := c.GetUserLikeCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *userLikeCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetUserLikeCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
