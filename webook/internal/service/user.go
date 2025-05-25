package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailDuplicate         = repository.ErrEmailDuplicate
	ErrInvalidEmailOrPassword = repository.ErrInvalidEmailOrPassword
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		repo: userRepository,
	}
}

func (u *UserService) Create(ctx context.Context, user domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encrypted)
	return u.repo.Create(ctx, user)
}

func (u *UserService) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return u.repo.FindByEmail(ctx, email)
}

func (u *UserService) Edit(ctx *gin.Context, userInfo domain.UserInfo) error {
	return u.repo.Edit(ctx, userInfo)
}

func (u *UserService) FindUserById(ctx *gin.Context, id int64) (domain.User, error) {
	return u.repo.FindUserById(ctx, id)
}

func (u *UserService) FindUserInfoById(ctx *gin.Context, id int64) (domain.UserInfo, error) {
	return u.repo.FindUserInfoById(ctx, id)
}
