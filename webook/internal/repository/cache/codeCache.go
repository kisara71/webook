package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone string, code int) error
	Verify(ctx context.Context, biz, phone string, code int) (bool, error)
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return newRedisCodeCache(client)
}

//go:embed lua/set_code.lua
var setCodeScript string

//go:embed lua/verify_code.lua
var varifyCodeScript string

var (
	ErrSendTooFrequent      = errors.New("send code too frequent")
	ErrSystemError          = errors.New("system error")
	ErrInvalidCode          = errors.New("invalid code")
	ErrTooManyVerifications = errors.New("too many verifications")
	ErrWrongCode            = errors.New("wrong code")
)

type redisCodeCache struct {
	client redis.Cmdable
}

func newRedisCodeCache(client redis.Cmdable) CodeCache {
	return &redisCodeCache{
		client: client,
	}
}

func (c *redisCodeCache) Set(ctx context.Context, biz, phone string, code int) error {
	res, err := c.client.Eval(ctx, setCodeScript, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -2:
		return ErrSendTooFrequent
	default:
		return ErrSystemError

	}
}
func (c *redisCodeCache) Verify(ctx context.Context, biz, phone string, code int) (bool, error) {
	res, err := c.client.Eval(ctx, varifyCodeScript, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return false, nil
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrTooManyVerifications
	case -2:
		return false, ErrWrongCode
	case -3:
		return false, ErrInvalidCode
	default:
		return false, ErrSystemError
	}
}
func (c *redisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
