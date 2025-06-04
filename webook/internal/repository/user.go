package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/cache/redisCache"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
)

var (
	ErrEmailDuplicate         = dao.ErrEmailDuplicate
	ErrInvalidEmailOrPassword = dao.ErrInvalidEmailOrPassword
	ErrRecordNotExist         = dao.ErrRecordNotFound
)

type UserRepository struct {
	userDao   *dao.UserDao
	userCache cache.UserCache
}

func NewUserRepository(userDao *dao.UserDao, userCache cache.UserCache) *UserRepository {
	return &UserRepository{
		userDao:   userDao,
		userCache: userCache,
	}
}

func (u *UserRepository) Create(ctx context.Context, user domain.User) error {
	return u.userDao.Insert(ctx, dao.UserPO{
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Password: user.Password,
	})
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return u.userDao.FindByEmail(ctx, email)
}

func (u *UserRepository) Edit(ctx context.Context, info domain.User) error {
	err := u.userDao.Edit(ctx, info)
	if err != nil {
		return err
	}
	go func() {
		newUser, _ := u.userDao.FindUserById(ctx, info.Id)
		err = u.userCache.Set(ctx, newUser)
		if err != nil {
			// log
		}
	}()
	return nil
}

func (u *UserRepository) FindUserById(ctx context.Context, id int64) (domain.User, error) {
	if user, err := u.userCache.Get(ctx, id); err == nil {
		return user, nil
	} else if errors.Is(err, redisCache.ErrKeyNotFound) {
		du, err := u.userDao.FindUserById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		go func() {
			newUser, _ := u.userDao.FindUserById(ctx, id)
			err = u.userCache.Set(ctx, newUser)
			if err != nil {
				// log
			}
		}()
		return du, nil
	}

	return domain.User{}, errors.New("redisCache error")

}
func (u *UserRepository) FindUser(ctx context.Context, filed string, value any) (domain.User, error) {
	return u.userDao.FindUser(ctx, filed, value)
}

func (u *UserRepository) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	ud, err := u.FindUser(ctx, "Phone", phone)
	if err == nil {
		return ud, nil
	} else if errors.Is(err, ErrRecordNotExist) {
		err = u.Create(ctx, domain.User{
			Phone: phone,
		})
		if err != nil {
			return ud, err
		}
		return u.FindUser(ctx, "Phone", phone)
	}
	return ud, err
}

//func (u *UserRepository) FindUserInfoById(ctx *gin.Context, id int64) (domain.User, error) {
//	if user, err := u.userCache.Get(ctx, id); err == nil {
//		return user, nil
//	}
//	return u.userDao.FindUserInfoById(ctx, id)
//}
