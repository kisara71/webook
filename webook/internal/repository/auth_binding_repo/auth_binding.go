package auth_binding_repo

import (
	"context"
	"errors"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
)

var (
	ErrSystemError = errors.New("系统错误")
)

type Repository interface {
	FindOrCreateOauth2Binding(ctx context.Context, binding domain.Oauth2Binding) (domain.Oauth2Binding, error)
}

type repositoryV1 struct {
	authDao dao.Dao
}

func NewRepository(authDao dao.Dao) Repository {
	return &repositoryV1{
		authDao: authDao,
	}
}

func (r *repositoryV1) FindOrCreateOauth2Binding(ctx context.Context, binding domain.Oauth2Binding) (domain.Oauth2Binding, error) {
	res, err := r.authDao.FindBinding(ctx, binding.Provider.ToString(), binding.ExternalID)
	if err != nil {
		if errors.Is(err, dao.ErrRecordNotFound) {
			res, err = r.authDao.InsertOauth2Binding(ctx, binding)
			if err != nil {
				return domain.Oauth2Binding{}, err
			}
			return res, nil
		} else {
			return domain.Oauth2Binding{}, ErrSystemError
		}
	}
	return res, nil
}
