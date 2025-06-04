package cache

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/domain"
)

type UserCache interface {
	Set(ctx context.Context, du domain.User) error
	Get(ctx context.Context, id int64) (domain.User, error)
}

type CodeCache interface {
	Set(ctx context.Context, biz, phone string, code int) error
	Verify(ctx context.Context, biz, phone string, code int) (bool, error)
}
