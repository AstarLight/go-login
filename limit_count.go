package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	ratelimiter "github.com/ilam01/limits-go"
	"time"
)

/*
	限流插件，支持单机模式和分布式模式，默认是redis模式

	单机模式基于进程内的内存存储，限流粒度是进程
	分布式模式基于redis实现
*/

const (
	LIMIT_COUNT_PRFIX = "limit-count-"
)

var (
	MinutelimiterPool map[string]*ratelimiter.Limiter
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

func Ratelimitinit() {
	MinutelimiterPool = make(map[string]*ratelimiter.Limiter)
}

func NewUserLimiter(username string, max int, mode string) *ratelimiter.Limiter {
	var limiter *ratelimiter.Limiter
	if mode == "redis" {
		limiter = ratelimiter.New(ratelimiter.Options{
			Max:      max,
			Duration: time.Hour,
			Client:   &redisClient{RedisDb},
			Ctx:      ctxRedis,
		})
	} else {
		limiter = ratelimiter.New(ratelimiter.Options{
			Max:      max,
			Duration: time.Hour,
		})
	}

	return limiter
}

// it to the upstream.
type LimitCount struct {

}

type LimitCountConf struct {
	Mode string `json:"mode"` // memory or redis
}

func (p *LimitCount) Name() string {
	return "pazhoulab-limit-count"
}

func IsReachLimit(string key, mode string) bool {
	if mode == "user" {
		username := key
		if _, ok := MinutelimiterPool[username]; !ok {
			newLimiter := NewUserLimiter(username, Conf.Common.LimitCount, "redis")
			MinutelimiterPool[username] = newLimiter

			skey := LIMIT_COUNT_PRFIX + username
		
			minutelimiter := MinutelimiterPool[username]

			res, err := minutelimiter.Get(ctxRedis, skey)
			if err != nil {
				// 取redis报错，认为是没有到达阈值
				return false
			}

			if res.Remaining <= 0 {
				// reach limit
				return true
			}
		}
	} else if  mode == "ip" {

	}

	return false
}
