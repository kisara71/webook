package oauth2

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/domain"
)

type Service interface {
	AuthURL(ctx context.Context) string
	ExchangeCode(ctx context.Context, code string, state string) (domain.Oauth2Binding, error)
}
