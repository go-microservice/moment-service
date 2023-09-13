package cache

//go:generate mockgen -source=internal/cache/comment_hot_cache.go -destination=internal/mock/comment_hot_cache_mock.go  -package mock

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/spf13/cast"

	"github.com/go-eagle/eagle/pkg/log"
	rdb "github.com/go-eagle/eagle/pkg/redis"
	redis "github.com/redis/go-redis/v9"

	"github.com/go-microservice/moment-service/internal/model"
)

const (
	// PrefixCommentHotCacheKey cache prefix
	PrefixCommentHotCacheKey = "comment:hot:%d"
)

// CommentHotCache define cache interface
type CommentHotCache interface {
	SetCommentHotCache(ctx context.Context, id int64, data *model.CommentHotModel, duration time.Duration) error
	GetCommentHotCache(ctx context.Context, id int64) (data *model.CommentHotModel, err error)
	GetListCommentHotCache(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error)
	DelCommentHotCache(ctx context.Context, id int64) error
}

// commentHotCache define cache struct
type commentHotCache struct {
	cache *redis.Client
}

// NewCommentHotCache new a cache
func NewCommentHotCache() CommentHotCache {
	return &commentHotCache{
		cache: rdb.RedisClient,
	}
}

// GetCommentHotCacheKey get cache key
func (c *commentHotCache) GetCommentHotCacheKey(id int64) string {
	return fmt.Sprintf(PrefixCommentHotCacheKey, id)
}

// SetCommentHotCache write to cache
func (c *commentHotCache) SetCommentHotCache(ctx context.Context, id int64, data *model.CommentHotModel, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCommentHotCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetCommentHotCache get from cache
func (c *commentHotCache) GetCommentHotCache(ctx context.Context, id int64) (data *model.CommentHotModel, err error) {
	cacheKey := c.GetCommentHotCacheKey(id)
	err = c.cache.Get(ctx, cacheKey).Err()
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

func (c *commentHotCache) GetListCommentHotCache(ctx context.Context, postID int64, lastID int64, limit int) (ret []int64, err error) {
	if postID == 0 {
		return nil, nil
	}

	cacheKey := c.GetCommentHotCacheKey(postID)
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

// DelCommentHotCache delete cache
func (c *commentHotCache) DelCommentHotCache(ctx context.Context, id int64) error {
	cacheKey := c.GetCommentHotCacheKey(id)
	err := c.cache.ZRem(ctx, cacheKey, id).Err()
	if err != nil {
		return err
	}
	return nil
}
