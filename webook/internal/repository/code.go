package repository

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
)

var (
	ErrSendTooFrequent      = cache.ErrSendTooFrequent
	ErrSystemError          = cache.ErrSystemError
	ErrInvalidCode          = cache.ErrInvalidCode
	ErrTooManyVerifications = cache.ErrTooManyVerifications
	ErrWrongCode            = cache.ErrWrongCode
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (c *CodeRepository) Store(ctx context.Context, biz, phone string, code int) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
