package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache interface {
	Set(ctx context.Context, du domain.User) error
	Get(ctx context.Context, id int64) (domain.User, error)
}

func NewUserCache(client redis.Cmdable) UserCache {
	return newRedisUserCache(client, time.Minute*15)
}

var (
	ErrKeyNotFound = redis.Nil
)

type redisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func newRedisUserCache(client redis.Cmdable, expiration time.Duration) UserCache {
	return &redisUserCache{
		cmd:        client,
		expiration: expiration,
	}
}

func (cache *redisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	var user domain.User
	val, err := cache.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	err = json.Unmarshal(val, &user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (cache *redisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *redisUserCache) Set(ctx context.Context, du domain.User) error {
	key := cache.key(du.Id)
	val, err := json.Marshal(&du)
	if err != nil {
		// log
		return errors.New("json marshal failed")
	}
	err = cache.cmd.Set(ctx, key, val, cache.expiration).Err()
	if err != nil {
		// log
		return err
	}
	return nil
}
