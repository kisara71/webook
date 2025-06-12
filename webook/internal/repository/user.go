package repository

import (
	"context"
	"database/sql"
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

func (U *userRepositoryV1) Create(ctx context.Context, user domain.User) error {
	return U.userDao.Insert(ctx, dao.UserEntity{
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

func (U *userRepositoryV1) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return U.userDao.FindByEmail(ctx, email)
}

func (U *userRepositoryV1) Edit(ctx context.Context, info domain.User) error {
	err := U.userDao.Edit(ctx, info)
	if err != nil {
		return err
	}
	go func() {
		newUser, _ := U.userDao.FindById(ctx, info.Id)
		err = U.userCache.Set(ctx, newUser)
		if err != nil {
			// log
		}
	}()
	return nil
}

func (U *userRepositoryV1) FindById(ctx context.Context, id int64) (domain.User, error) {
	if user, err := U.userCache.Get(ctx, id); err == nil {
		return user, nil
	} else if errors.Is(err, cache.ErrKeyNotFound) {
		du, err := U.userDao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		//go func() {
		err = U.userCache.Set(ctx, du)
		if err != nil {
			// log
		}
		//}()
		return du, nil
	}

	return domain.User{}, errors.New("redisCache error")

}
func (U *userRepositoryV1) FindUser(ctx context.Context, filed string, value any) (domain.User, error) {
	return U.userDao.FindUser(ctx, filed, value)
}

func (U *userRepositoryV1) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	ud, err := U.FindUser(ctx, "Phone", phone)
	if err == nil {
		return ud, nil
	} else if errors.Is(err, dao.ErrRecordNotFound) {
		err = U.Create(ctx, domain.User{
			Phone: phone,
		})
		if err != nil {
			return ud, err
		}
		return U.FindUser(ctx, "Phone", phone)
	}
	return ud, err
}

//func (u *userRepositoryV1) FindUserInfoById(ctx *gin.Context, id int64) (domain.User, error) {
//	if user, err := u.userCache.Get(ctx, id); err == nil {
//		return user, nil
//	}
//	return u.userDao.FindUserInfoById(ctx, id)
//}
