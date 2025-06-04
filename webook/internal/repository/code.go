package repository

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/cache/redisCache"
)

var (
	ErrSendTooFrequent      = redisCache.ErrSendTooFrequent
	ErrSystemError          = redisCache.ErrSystemError
	ErrInvalidCode          = redisCache.ErrInvalidCode
	ErrTooManyVerifications = redisCache.ErrTooManyVerifications
	ErrWrongCode            = redisCache.ErrWrongCode
)

type CodeRepository struct {
	codeCache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		codeCache: c,
	}
}

func (c *CodeRepository) Store(ctx context.Context, biz, phone string, code int) error {
	return c.codeCache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error) {
	return c.codeCache.Verify(ctx, biz, phone, code)
}
