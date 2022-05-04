package cache

//go:generate mockgen -source=internal/cache/comment_index_cache.go -destination=internal/mock/comment_index_cache_mock.go  -package mock

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/spf13/cast"

	rdb "github.com/go-eagle/eagle/pkg/redis"
	"github.com/go-redis/redis/v8"
)

const (
	// PrefixCommentIndexCacheKey cache prefix oid_type_sort
	PrefixCommentIndexCacheKey = "comment:index:%d_%d_%d"
)

// CommentIndexCache define cache interface
type CommentIndexCache interface {
	SetCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, commentID int64, score float64, duration time.Duration) error
	GetListCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, lastID int64, limit int) (ret []int64, err error)
	MultiSetCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, cmtIDs []int64, scores []float64, duration time.Duration) error
	DelCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, cmtID int64) error
}

// commentIndexCache define cache struct
type commentIndexCache struct {
	cache *redis.Client
}

// NewCommentIndexCache new a cache
func NewCommentIndexCache() CommentIndexCache {
	return &commentIndexCache{
		cache: rdb.RedisClient,
	}
}

// GetCommentIndexCacheKey get cache key
func (c *commentIndexCache) GetCommentIndexCacheKey(objID int64, objType int, sortType int) string {
	return fmt.Sprintf(PrefixCommentIndexCacheKey, objID, objType, sortType)
}

// SetCommentIndexCache write to cache
func (c *commentIndexCache) SetCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, commentID int64, score float64, duration time.Duration) error {
	if objID == 0 || commentID == 0 {
		return nil
	}

	cacheKey := c.GetCommentIndexCacheKey(objID, objType, sortType)
	pipe := c.cache.Pipeline()
	err := pipe.ZAdd(ctx, cacheKey, &redis.Z{
		Score:  score,
		Member: commentID,
	}).Err()
	if err != nil {
		return err
	}
	err = pipe.Expire(ctx, cacheKey, duration).Err()
	if err != nil {
		return err
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetCommentIndexCache get list from cache
func (c *commentIndexCache) GetListCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, lastID int64, limit int) (ret []int64, err error) {
	if objID == 0 {
		return nil, nil
	}

	cacheKey := c.GetCommentIndexCacheKey(objID, objType, sortType)
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

// MultiSetCommentIndexCache batch set cache
func (c *commentIndexCache) MultiSetCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, cmtIDs []int64, scores []float64, duration time.Duration) error {
	if objID == 0 {
		return errors.New("objID is empty")
	}
	if len(cmtIDs) == 0 {
		return errors.New("cmdIDs is empty")
	}
	if len(scores) == 0 {
		return errors.New("scores is empty")
	}
	if len(cmtIDs) != len(scores) {
		return errors.New("cmdIDs and score length not equal")
	}

	cacheKey := c.GetCommentIndexCacheKey(objID, objType, sortType)
	var data []*redis.Z
	for k, v := range cmtIDs {
		data = append(data, &redis.Z{
			Score:  scores[k],
			Member: v,
		})
	}
	err := c.cache.ZAdd(ctx, cacheKey, data...).Err()
	if err != nil {
		return err
	}

	return nil
}

// DelCommentIndexCache delete cache
func (c *commentIndexCache) DelCommentIndexCache(ctx context.Context, objID int64, objType int, sortType int, cmtID int64) error {
	cacheKey := c.GetCommentIndexCacheKey(objID, objType, sortType)
	err := c.cache.ZRem(ctx, cacheKey, cmtID).Err()
	if err != nil {
		return err
	}
	return nil
}
