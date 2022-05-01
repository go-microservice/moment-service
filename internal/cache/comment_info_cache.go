package cache

//go:generate mockgen -source=internal/cache/comment_info_cache.go -destination=internal/mock/comment_info_cache_mock.go  -package mock

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
	// PrefixCommentInfoCacheKey cache prefix
	PrefixCommentInfoCacheKey = "comment:info:%d"
)

// CommentInfoCache define cache interface
type CommentInfoCache interface {
	SetCommentInfoCache(ctx context.Context, id int64, data *model.CommentInfoModel, duration time.Duration) error
	GetCommentInfoCache(ctx context.Context, id int64) (data *model.CommentInfoModel, err error)
	MultiGetCommentInfoCache(ctx context.Context, ids []int64) (map[string]*model.CommentInfoModel, error)
	MultiSetCommentInfoCache(ctx context.Context, data []*model.CommentInfoModel, duration time.Duration) error
	DelCommentInfoCache(ctx context.Context, id int64) error
}

// commentInfoCache define cache struct
type commentInfoCache struct {
	cache cache.Cache
}

// NewCommentInfoCache new a cache
func NewCommentInfoCache() CommentInfoCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &commentInfoCache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return &model.CommentInfoModel{}
		}),
	}
}

// GetCommentInfoCacheKey get cache key
func (c *commentInfoCache) GetCommentInfoCacheKey(id int64) string {
	return fmt.Sprintf(PrefixCommentInfoCacheKey, id)
}

// SetCommentInfoCache write to cache
func (c *commentInfoCache) SetCommentInfoCache(ctx context.Context, id int64, data *model.CommentInfoModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCommentInfoCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetCommentInfoCache get from cache
func (c *commentInfoCache) GetCommentInfoCache(ctx context.Context, id int64) (data *model.CommentInfoModel, err error) {
	cacheKey := c.GetCommentInfoCacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetCommentInfoCache batch get cache
func (c *commentInfoCache) MultiGetCommentInfoCache(ctx context.Context, ids []int64) (map[string]*model.CommentInfoModel, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCommentInfoCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.CommentInfoModel)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSetCommentInfoCache batch set cache
func (c *commentInfoCache) MultiSetCommentInfoCache(ctx context.Context, data []*model.CommentInfoModel, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCommentInfoCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// DelCommentInfoCache delete cache
func (c *commentInfoCache) DelCommentInfoCache(ctx context.Context, id int64) error {
	cacheKey := c.GetCommentInfoCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *commentInfoCache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.GetCommentInfoCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
