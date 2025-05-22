package repository

import (
	"context"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
)

type UserRepository struct {
	userDao *dao.UserDao
}

func NewUserRepository(userDao *dao.UserDao) *UserRepository {
	return &UserRepository{
		userDao: userDao,
	}
}

func (u *UserRepository) Create(ctx context.Context, user domain.User) error {
	return u.userDao.Insert(ctx, dao.UserPO{
		Email:    user.Email,
		Password: user.Password,
	})
}
