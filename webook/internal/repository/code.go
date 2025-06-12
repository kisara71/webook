package repository

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone string, code int) error
	VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error)
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return newCodeRepositoryV1(c)
}

var (
	ErrSendTooFrequent      = cache.ErrSendTooFrequent
	ErrSystemError          = cache.ErrSystemError
	ErrInvalidCode          = cache.ErrInvalidCode
	ErrTooManyVerifications = cache.ErrTooManyVerifications
	ErrWrongCode            = cache.ErrWrongCode
)

type codeRepositoryV1 struct {
	codeCache cache.CodeCache
}

func newCodeRepositoryV1(c cache.CodeCache) CodeRepository {
	return &codeRepositoryV1{
		codeCache: c,
	}
}

func (R *codeRepositoryV1) Store(ctx context.Context, biz, phone string, code int) error {
	return R.codeCache.Set(ctx, biz, phone, code)
}

func (R *codeRepositoryV1) VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error) {
	return R.codeCache.Verify(ctx, biz, phone, code)
}
