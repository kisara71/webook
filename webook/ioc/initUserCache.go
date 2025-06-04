package ioc

import (
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/cache/redisCache"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitUserCache(client redis.Cmdable) cache.UserCache {
	return redisCache.NewUserCache(client, time.Minute*15)
}
