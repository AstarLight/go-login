package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	ratelimiter "github.com/ilam01/limits-go"
	"time"
)

/*
	限流插件，支持单机模式和分布式模式，默认是redis模式

	单机模式基于进程内的内存存储，限流粒度是进程
	分布式模式基于redis实现
*/

var (
	userlimiterPool map[string]*ratelimiter.Limiter
	iplimiterPool   map[string]*ratelimiter.Limiter

	ratelimitCtx   = context.Background()
	ratelimitRedis *redis.Client
)

// Implements RedisClient for redis.Client
type redisClient struct {
	*redis.Client
}

func (c *redisClient) RateDel(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

func (c *redisClient) RateEvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return c.EvalSha(ctx, sha1, keys, args...).Result()
}

func (c *redisClient) RateScriptLoad(ctx context.Context, script string) (string, error) {
	return c.ScriptLoad(ctx, script).Result()
}

func RatelimitInit() {
	userlimiterPool = make(map[string]*ratelimiter.Limiter)
	iplimiterPool = make(map[string]*ratelimiter.Limiter)

	ratelimitRedis = redis.NewClient(&redis.Options{
		Addr: Conf.Redis.Addr,
	})

	pong, err := RedisDb.Ping(ratelimitCtx).Result()
	if err != nil {
		panic(fmt.Sprintf("connect redis fail: %v", err))
	} else {
		fmt.Println("connect redis succ,", pong)
	}
}

func NewLimiter(max int, mode string) *ratelimiter.Limiter {
	var limiter *ratelimiter.Limiter
	if mode == "redis" {
		limiter = ratelimiter.New(ratelimiter.Options{
			Max:      max,
			Duration: time.Minute,
			Client:   &redisClient{ratelimitRedis},
			Ctx:      ratelimitCtx,
		})
	} else {
		limiter = ratelimiter.New(ratelimiter.Options{
			Max:      max,
			Duration: time.Minute,
		})
	}

	return limiter
}

func GetKey(key string) string {
	return Conf.RateLimit.KeyPrefix + key
}

func IsReachLimit(key, mode string) bool {
	var limiterPool map[string]*ratelimiter.Limiter
	if mode == "user" {
		limiterPool = userlimiterPool
	} else if mode == "ip" {
		limiterPool = iplimiterPool
	} else {
		fmt.Println("error config ratelimit mode!")
		return false
	}

	if _, ok := limiterPool[key]; !ok {
		newLimiter := NewLimiter(Conf.RateLimit.Limit, Conf.RateLimit.Mode)
		limiterPool[key] = newLimiter
	}

	limiter := limiterPool[key]

	rkey := GetKey(key)
	res, err := limiter.Get(ratelimitCtx, rkey)
	if err != nil {
		// 取redis报错，认为是没有到达阈值
		return false
	}

	if res.Remaining <= 0 {
		// reach limit
		return true
	}

	return false
}
