package auth_binding_service

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/auth_binding_repo"
)

type Service interface {
	FindOrCreateOauth2Binding(ctx context.Context, binding domain.Oauth2Binding) (domain.Oauth2Binding, error)
}

type serviceV1 struct {
	repo auth_binding_repo.Repository
}

func NewService(repo auth_binding_repo.Repository) Service {
	return &serviceV1{
		repo: repo,
	}
}

func (s *serviceV1) FindOrCreateOauth2Binding(ctx context.Context, binding domain.Oauth2Binding) (domain.Oauth2Binding, error) {
	return s.repo.FindOrCreateOauth2Binding(ctx, binding)
}
