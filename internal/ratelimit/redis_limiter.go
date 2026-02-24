package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(addr string) *RedisLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisLimiter{client: rdb}
}

func (rl *RedisLimiter) IsAllowed(ctx context.Context, userID string, limit int) (bool, error) {
	now := time.Now().UnixNano()
	window := time.Minute.Nanoseconds()
	minimum := now - window

	key := "rate_limit:" + userID
	//using Pipeline for increase speed and Atomic
	pipe := rl.client.TxPipeline()

	//1 remove old requests
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(minimum, 10))

	//2 count remained requests in window
	count := pipe.ZCard(ctx, key)

	//3 add current requets
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})

	//4 set expire for auto clean
	pipe.Expire(ctx, key, time.Minute)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	//check is request reached to limit or no
	return int(count.Val()) < limit, nil
}
