package cache

//go:generate mockgen -source=internal/cache/comment_latest_cache.go -destination=internal/mock/comment_latest_cache_mock.go  -package mock

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/spf13/cast"

	"github.com/go-eagle/eagle/pkg/log"
	rdb "github.com/go-eagle/eagle/pkg/redis"
	"github.com/go-redis/redis/v8"

	"github.com/go-microservice/moment-service/internal/model"
)

const (
	// PrefixCommentLatestCacheKey cache prefix
	PrefixCommentLatestCacheKey = "comment:latest:%d"
)

// CommentLatestCache define cache interface
type CommentLatestCache interface {
	SetCommentLatestCache(ctx context.Context, id int64, data *model.CommentLatestModel, duration time.Duration) error
	GetCommentLatestCache(ctx context.Context, id int64) (data *model.CommentLatestModel, err error)
	GetListCommentLatestCache(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error)
	DelCommentLatestCache(ctx context.Context, id int64) error
}

// commentLatestCache define cache struct
type commentLatestCache struct {
	cache *redis.Client
}

// NewCommentLatestCache new a cache
func NewCommentLatestCache() CommentLatestCache {
	return &commentLatestCache{
		cache: rdb.RedisClient,
	}
}

// GetCommentLatestCacheKey get cache key
func (c *commentLatestCache) GetCommentLatestCacheKey(id int64) string {
	return fmt.Sprintf(PrefixCommentLatestCacheKey, id)
}

// SetCommentLatestCache write to cache
func (c *commentLatestCache) SetCommentLatestCache(ctx context.Context, id int64, data *model.CommentLatestModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCommentLatestCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetCommentLatestCache get from cache
func (c *commentLatestCache) GetCommentLatestCache(ctx context.Context, id int64) (data *model.CommentLatestModel, err error) {
	cacheKey := c.GetCommentLatestCacheKey(id)
	err = c.cache.Get(ctx, cacheKey).Err()
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

func (c *commentLatestCache) GetListCommentLatestCache(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error) {
	if postID == 0 {
		return nil, nil
	}

	cacheKey := c.GetCommentLatestCacheKey(postID)
	// get score by lastID
	score := "+inf"
	if lastID != int64(math.MaxInt64) {
		s, err := c.cache.ZScore(ctx, cacheKey, cast.ToString(lastID)).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		score = cast.ToString(s)
	}

	// get data by score
	data, err := c.cache.ZRevRangeByScore(ctx, cacheKey, &redis.ZRangeBy{
		Min:    score,
		Max:    "-inf",
		Offset: 0,
		Count:  int64(limit),
	}).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	for _, v := range data {
		ret = append(ret, cast.ToInt64(v))
	}

	return ret, nil
}

// DelCommentLatestCache delete cache
func (c *commentLatestCache) DelCommentLatestCache(ctx context.Context, id int64) error {
	cacheKey := c.GetCommentLatestCacheKey(id)
	err := c.cache.ZRem(ctx, cacheKey, id).Err()
	if err != nil {
		return err
	}
	return nil
}
