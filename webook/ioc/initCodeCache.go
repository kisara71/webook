package ioc

import (
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/cache/redisCache"
	"github.com/redis/go-redis/v9"
)

func InitCodeCache(client redis.Cmdable) cache.CodeCache {
	return redisCache.NewCodeCache(client)
}
