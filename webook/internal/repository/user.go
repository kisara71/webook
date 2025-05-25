package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
)

var (
	ErrEmailDuplicate         = dao.ErrEmailDuplicate
	ErrInvalidEmailOrPassword = dao.ErrInvalidEmailOrPassword
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

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return u.userDao.FindByEmail(ctx, email)
}

func (u *UserRepository) Edit(ctx *gin.Context, info domain.UserInfo) error {
	return u.userDao.Edit(ctx, info)
}

func (u *UserRepository) FindUserById(ctx *gin.Context, id int64) (domain.User, error) {
	return u.userDao.FindUserById(ctx, id)
}

func (u *UserRepository) FindUserInfoById(ctx *gin.Context, id int64) (domain.UserInfo, error) {
	return u.userDao.FindUserInfoById(ctx, id)
}
