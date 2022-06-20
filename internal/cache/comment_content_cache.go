package cache

//go:generate mockgen -source=internal/cache/comment_content_cache.go -destination=internal/mock/comment_content_cache_mock.go  -package mock

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
	// PrefixCommentContentCacheKey cache prefix
	PrefixCommentContentCacheKey = "comment:content:%d"
)

// CommentContentCache define cache interface
type CommentContentCache interface {
	SetCommentContentCache(ctx context.Context, id int64, data *model.CommentContentModel, duration time.Duration) error
	GetCommentContentCache(ctx context.Context, id int64) (data *model.CommentContentModel, err error)
	MultiGetCommentContentCache(ctx context.Context, ids []int64) (map[int64]*model.CommentContentModel, error)
	MultiSetCommentContentCache(ctx context.Context, data []*model.CommentContentModel, duration time.Duration) error
	DelCommentContentCache(ctx context.Context, id int64) error
}

// commentContentCache define cache struct
type commentContentCache struct {
	cache cache.Cache
}

// NewCommentContentCache new a cache
func NewCommentContentCache() CommentContentCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &commentContentCache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return &model.CommentContentModel{}
		}),
	}
}

// GetCommentContentCacheKey get cache key
func (c *commentContentCache) GetCommentContentCacheKey(id int64) string {
	return fmt.Sprintf(PrefixCommentContentCacheKey, id)
}

// SetCommentContentCache write to cache
func (c *commentContentCache) SetCommentContentCache(ctx context.Context, id int64, data *model.CommentContentModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCommentContentCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetCommentContentCache get from cache
func (c *commentContentCache) GetCommentContentCache(ctx context.Context, id int64) (data *model.CommentContentModel, err error) {
	cacheKey := c.GetCommentContentCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil && err != redis.ErrRedisNotFound {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetCommentContentCache batch get cache
func (c *commentContentCache) MultiGetCommentContentCache(ctx context.Context, ids []int64) (map[int64]*model.CommentContentModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCommentContentCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	cacheMap := make(map[string]*model.CommentContentModel)
	err := c.cache.MultiGet(ctx, keys, cacheMap)
	if err != nil {
		return nil, err
	}
	retMap := make(map[int64]*model.CommentContentModel)
	for _, v := range ids {
		cacheKey := c.GetCommentContentCacheKey(v)
		val, ok := cacheMap[cacheKey]
		if ok {
			retMap[v] = val
		}
	}
	return retMap, nil
}

// MultiSetCommentContentCache batch set cache
func (c *commentContentCache) MultiSetCommentContentCache(ctx context.Context, data []*model.CommentContentModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCommentContentCacheKey(v.CommentId)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelCommentContentCache delete cache
func (c *commentContentCache) DelCommentContentCache(ctx context.Context, id int64) error {
	cacheKey := c.GetCommentContentCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *commentContentCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetCommentContentCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
