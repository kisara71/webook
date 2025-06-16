package service

import (
	"context"
	"errors"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	Edit(ctx context.Context, userInfo domain.User) error
	FindUserByEmail(ctx context.Context, email string) (domain.User, error)
	FindUserById(ctx context.Context, id int64) (domain.User, error)
	FindUser(ctx context.Context, filed string, value any) (domain.User, error)
	FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error)
	Login(ctx context.Context, email string, password string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

func NewUserService(up repository.UserRepository) UserService {
	return newUserServiceV1(up)
}

var (
	ErrEmailDuplicate         = repository.ErrEmailDuplicate
	ErrInvalidEmailOrPassword = repository.ErrInvalidEmailOrPassword
	ErrUserNotExist           = repository.ErrRecordNotExist
)

type userServiceV1 struct {
	repo repository.UserRepository
}

func (u *userServiceV1) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	return u.repo.FindOrCreateByOpenID(ctx, info.OpenID)
}

func newUserServiceV1(userRepository repository.UserRepository) UserService {
	return &userServiceV1{
		repo: userRepository,
	}
}

func (u *userServiceV1) SignUp(ctx context.Context, user domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encrypted)
	return u.repo.Create(ctx, user)
}

func (u *userServiceV1) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return u.repo.FindByEmail(ctx, email)
}

func (u *userServiceV1) Edit(ctx context.Context, userInfo domain.User) error {
	return u.repo.Edit(ctx, userInfo)
}

func (u *userServiceV1) FindUserById(ctx context.Context, id int64) (domain.User, error) {
	return u.repo.FindById(ctx, id)
}
func (u *userServiceV1) FindUser(ctx context.Context, filed string, value any) (domain.User, error) {
	return u.repo.FindUser(ctx, filed, value)
}
func (u *userServiceV1) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	return u.repo.FindOrCreateByPhone(ctx, phone)
}
func (u *userServiceV1) Login(ctx context.Context, email string, password string) (domain.User, error) {
	du, err := u.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotExist) {
			return domain.User{}, ErrUserNotExist
		}
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(du.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidEmailOrPassword
	}
	return du, nil
}

//func (u *userServiceV1) FindUserInfoById(ctx *gin.Context, id int64) (domain.User, error) {
//	return u.repo.FindUserInfoById(ctx, id)
//}
