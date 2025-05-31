package repository

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
)

var (
	ErrEmailDuplicate         = dao.ErrEmailDuplicate
	ErrInvalidEmailOrPassword = dao.ErrInvalidEmailOrPassword
)

type UserRepository struct {
	userDao   *dao.UserDao
	userCache *cache.UserCache
}

func NewUserRepository(userDao *dao.UserDao, userCache *cache.UserCache) *UserRepository {
	return &UserRepository{
		userDao:   userDao,
		userCache: userCache,
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

func (u *UserRepository) Edit(ctx *gin.Context, info domain.User) error {
	err := u.userDao.Edit(ctx, info)
	if err != nil {
		return err
	}
	newUser, _ := u.userDao.FindUserById(ctx, info.Id)
	err = u.userCache.Set(ctx, newUser)
	if err != nil {
		// log
	}
	return nil
}

func (u *UserRepository) FindUserById(ctx *gin.Context, id int64) (domain.User, error) {
	if user, err := u.userCache.Get(ctx, id); err == nil {
		return user, nil
	} else if errors.Is(err, cache.ErrKeyNotFound) {
		du, err := u.userDao.FindUserById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		err = u.userCache.Set(ctx, du)
		return du, nil
	}
	return domain.User{}, errors.New("redis error")

}

//func (u *UserRepository) FindUserInfoById(ctx *gin.Context, id int64) (domain.User, error) {
//	if user, err := u.userCache.Get(ctx, id); err == nil {
//		return user, nil
//	}
//	return u.userDao.FindUserInfoById(ctx, id)
//}
