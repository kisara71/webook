package repository

import (
	"context"
	"errors"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Edit(ctx context.Context, info domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error)
	FindUser(ctx context.Context, filed string, value any) (domain.User, error)
	FindOrCreateByOpenID(ctx context.Context, openID string) (domain.User, error)
}

func NewUserRepository(d dao.Dao, c cache.UserCache) UserRepository {
	return newUserRepositoryV1(d, c)
}

var (
	ErrEmailDuplicate         = dao.ErrEmailDuplicate
	ErrInvalidEmailOrPassword = dao.ErrInvalidEmailOrPassword
	ErrRecordNotExist         = dao.ErrRecordNotFound
)

type userRepositoryV1 struct {
	userDao   dao.Dao
	userCache cache.UserCache
}

func newUserRepositoryV1(userDao dao.Dao, userCache cache.UserCache) UserRepository {
	return &userRepositoryV1{
		userDao:   userDao,
		userCache: userCache,
	}
}
func (u *userRepositoryV1) FindOrCreateByOpenID(ctx context.Context, openID string) (domain.User, error) {
	user, err := u.userDao.FindUser(ctx, "OpenID", openID)
	if err == nil {
		return user, nil
	} else if errors.Is(err, dao.ErrRecordNotFound) {
		err = u.userDao.Insert(ctx, domain.User{
			WechatInfo: domain.WechatInfo{
				OpenID: openID,
			},
		})
		if err != nil {
			return domain.User{}, err
		}
		return u.FindUser(ctx, "OpenID", openID)
	}
	return domain.User{}, err
}
func (u *userRepositoryV1) Create(ctx context.Context, user domain.User) error {
	return u.userDao.Insert(ctx, user)
}

func (u *userRepositoryV1) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return u.userDao.FindByEmail(ctx, email)
}

func (u *userRepositoryV1) Edit(ctx context.Context, info domain.User) error {
	err := u.userDao.Edit(ctx, info)
	if err != nil {
		return err
	}
	go func() {
		newUser, _ := u.userDao.FindById(ctx, info.Id)
		err = u.userCache.Set(ctx, newUser)
		if err != nil {
			// log
		}
	}()
	return nil
}

func (u *userRepositoryV1) FindById(ctx context.Context, id int64) (domain.User, error) {
	if user, err := u.userCache.Get(ctx, id); err == nil {
		return user, nil
	} else if errors.Is(err, cache.ErrKeyNotFound) {
		du, err := u.userDao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		//go func() {
		err = u.userCache.Set(ctx, du)
		if err != nil {
			// log
		}
		//}()
		return du, nil
	}

	return domain.User{}, errors.New("redisCache error")

}
func (u *userRepositoryV1) FindUser(ctx context.Context, filed string, value any) (domain.User, error) {
	return u.userDao.FindUser(ctx, filed, value)
}

func (u *userRepositoryV1) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	ud, err := u.FindUser(ctx, "Phone", phone)
	if err == nil {
		return ud, nil
	} else if errors.Is(err, dao.ErrRecordNotFound) {
		err = u.Create(ctx, domain.User{
			Phone: phone,
		})
		if err != nil {
			return domain.User{}, err
		}
		return u.FindUser(ctx, "Phone", phone)
	}
	return domain.User{}, err
}

//func (u *userRepositoryV1) FindUserInfoById(ctx *gin.Context, id int64) (domain.User, error) {
//	if user, err := u.userCache.Get(ctx, id); err == nil {
//		return user, nil
//	}
//	return u.userDao.FindUserInfoById(ctx, id)
//}
