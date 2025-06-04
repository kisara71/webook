package redisCache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	ErrKeyNotFound = redis.Nil
)

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable, expiration time.Duration) *UserCache {
	return &UserCache{
		cmd:        client,
		expiration: expiration,
	}
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
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

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *UserCache) Set(ctx context.Context, du domain.User) error {
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
