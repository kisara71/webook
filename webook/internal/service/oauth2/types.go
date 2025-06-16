package oauth2

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/domain"
)

type Config struct {
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	AuthURLPattern    string `yaml:"auth_url_pattern"`
	GetBindURLPattern string `yaml:"get_bind_url_pattern"`
}

type Service interface {
	AuthURL(ctx context.Context) string
	ExchangeCode(ctx context.Context, code string, state string) (domain.Oauth2Binding, error)
}
