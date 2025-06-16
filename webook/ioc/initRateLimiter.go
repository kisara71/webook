package ioc

import (
	"github.com/kisara71/WeBook/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func initRateLimiter(cmd redis.Cmdable) ratelimit.RateLimiter {
	return ratelimit.NewRedisSlideWindowRateLimiter(cmd, time.Second, 3000)
}
