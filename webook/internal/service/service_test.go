package service

import (
	"context"
	"errors"
	"github.com/kisara71/WeBook/webook/internal/domain"
	repomocks "github.com/kisara71/WeBook/webook/internal/repository/mocks"
	"github.com/kisara71/WeBook/webook/internal/repository/user_repo"
	"github.com/kisara71/WeBook/webook/internal/service/user_service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserServiceV1_Login(t *testing.T) {

	testCases := []struct {
		name     string
		mocks    func(ctrl *gomock.Controller) user_repo.UserRepository
		user     domain.User
		wantErr  error
		wantUser domain.User
	}{
		{
			name: "登录成功",
			mocks: func(ctrl *gomock.Controller) user_repo.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$.hsnJY5Q4G0bUozFbKbrwe34Li8gqzAuVz8uZAN8SLC3TRtZBKiz2",
				}, nil)
				return repo
			},
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$.hsnJY5Q4G0bUozFbKbrwe34Li8gqzAuVz8uZAN8SLC3TRtZBKiz2",
			},
			wantErr: nil,
			user: domain.User{
				Email:    "123@qq.com",
				Password: "123456Q.",
			},
		},
		{
			name: "用户不存在",
			mocks: func(ctrl *gomock.Controller) user_repo.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, user_repo.ErrRecordNotExist)
				return repo
			},
			wantUser: domain.User{},
			wantErr:  user_service.ErrUserNotExist,
			user: domain.User{
				Email:    "123@qq.com",
				Password: "123456Q.",
			},
		},
		{
			name: "邮箱或密码错误",
			mocks: func(ctrl *gomock.Controller) user_repo.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$.hsnJY5Q4G0bUozFbKbrwe34Li8gqzAuVz8uZAN8SLC3TRtZBKiz2"}, nil)
				return repo
			},
			wantUser: domain.User{},
			wantErr:  user_service.ErrInvalidEmailOrPassword,
			user: domain.User{
				Email:    "123@qq.com",
				Password: "123456Q",
			},
		},
		{
			name: "系统错误",
			mocks: func(ctrl *gomock.Controller) user_repo.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("系统错误"))
				return repo
			},
			wantUser: domain.User{},
			wantErr:  errors.New("系统错误"),
			user: domain.User{
				Email:    "123@qq.com",
				Password: "123456Q",
			},
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usvc := user_service.NewUserService(ts.mocks(ctrl))
			user, err := usvc.Login(context.Background(), ts.user.Email, ts.user.Password)
			assert.Equal(t, err, ts.wantErr)
			assert.Equal(t, user, ts.wantUser)
		})
	}

}
func TestBcrypt(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456Q."), bcrypt.DefaultCost)
	t.Log(string(hashed))
}
