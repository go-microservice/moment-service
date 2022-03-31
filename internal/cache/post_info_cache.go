package cache

//go:generate mockgen -source=internal/cache/post_info_cache.go -destination=internal/mock/post_info_cache_mock.go  -package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"

	"github.com/go-microservice/moment-service/internal/model"
)

const (
	// PrefixPostInfoCacheKey cache prefix
	PrefixPostInfoCacheKey = "post:info:%d"
)

// PostInfoCache define cache interface
type PostInfoCache interface {
	SetPostInfoCache(ctx context.Context, id int64, data *model.PostInfoModel, duration time.Duration) error
	GetPostInfoCache(ctx context.Context, id int64) (data *model.PostInfoModel, err error)
	MultiGetPostInfoCache(ctx context.Context, ids []int64) (map[string]*model.PostInfoModel, error)
	MultiSetPostInfoCache(ctx context.Context, data []*model.PostInfoModel, duration time.Duration) error
	DelPostInfoCache(ctx context.Context, id int64) error
}

// postInfoCache define cache struct
type postInfoCache struct {
	cache cache.Cache
}

// NewPostInfoCache new a cache
func NewPostInfoCache() PostInfoCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &postInfoCache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return &model.PostInfoModel{}
		}),
	}
}

// GetPostInfoCacheKey get cache key
func (c *postInfoCache) GetPostInfoCacheKey(id int64) string {
	return fmt.Sprintf(PrefixPostInfoCacheKey, id)
}

// SetPostInfoCache write to cache
func (c *postInfoCache) SetPostInfoCache(ctx context.Context, id int64, data *model.PostInfoModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetPostInfoCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetPostInfoCache get from cache
func (c *postInfoCache) GetPostInfoCache(ctx context.Context, id int64) (data *model.PostInfoModel, err error) {
	cacheKey := c.GetPostInfoCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetPostInfoCache batch get cache
func (c *postInfoCache) MultiGetPostInfoCache(ctx context.Context, ids []int64) (map[string]*model.PostInfoModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetPostInfoCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.PostInfoModel)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSetPostInfoCache batch set cache
func (c *postInfoCache) MultiSetPostInfoCache(ctx context.Context, data []*model.PostInfoModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetPostInfoCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelPostInfoCache delete cache
func (c *postInfoCache) DelPostInfoCache(ctx context.Context, id int64) error {
	cacheKey := c.GetPostInfoCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *postInfoCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetPostInfoCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
