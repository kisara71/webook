package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/service"
	svcmock "github.com/kisara71/WeBook/webook/internal/service/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {

	testCases := []struct {
		name    string
		reqBody string
		getSvc  func(ctrl *gomock.Controller) service.UserService

		wantBody string
		wantCode int
	}{
		{
			name: "注册成功",
			reqBody: `
			{
				"email": "123@qq.com",
				"password": "123456Q.",
				"confirmPassword": "123456Q."
			}
			`,
			getSvc: func(ctrl *gomock.Controller) service.UserService {
				usvc := svcmock.NewMockUserService(ctrl)
				usvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "123456Q.",
				}).Return(nil)
				return usvc
			},
			wantBody: "注册成功",
			wantCode: 200,
		},
		{
			name:    "ctx Bind 失败",
			reqBody: `dfsfds`,
			getSvc: func(ctrl *gomock.Controller) service.UserService {
				return svcmock.NewMockUserService(ctrl)
			},
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
		{
			name: "密码不同",
			reqBody: `
			{
				"email": "123@qq.com",
				"password": "123456Q",
				"confirmPassword": "123456Q."
			}
			`,
			getSvc: func(ctrl *gomock.Controller) service.UserService {
				usvc := svcmock.NewMockUserService(ctrl)

				return usvc
			},
			wantBody: "两次密码不同",
			wantCode: 200,
		},
		{
			name: "邮箱格式错误",
			reqBody: `
			{
				"email": "123com",
				"password": "123456Q.",
				"confirmPassword": "123456Q."
			}
			`,
			getSvc: func(ctrl *gomock.Controller) service.UserService {
				usvc := svcmock.NewMockUserService(ctrl)
				return usvc
			},
			wantBody: "错误的邮箱格式",
			wantCode: 200,
		},
		{
			name: "密码格式错误",
			reqBody: `
			{
				"email": "123@qq.com",
				"password": "123456Q",
				"confirmPassword": "123456Q"
			}
			`,
			getSvc: func(ctrl *gomock.Controller) service.UserService {
				usvc := svcmock.NewMockUserService(ctrl)
				return usvc
			},
			wantBody: "错误的密码格式",
			wantCode: 200,
		},
		{
			name: "邮箱已注册",
			reqBody: `
			{
				"email": "123@qq.com",
				"password": "123456Q.",
				"confirmPassword": "123456Q."
			}
			`,
			getSvc: func(ctrl *gomock.Controller) service.UserService {
				usvc := svcmock.NewMockUserService(ctrl)
				usvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "123456Q.",
				}).Return(service.ErrEmailDuplicate)
				return usvc
			},
			wantBody: "邮箱已注册",
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			req, err := http.NewRequest(
				http.MethodPost,
				"/users/signup",
				bytes.NewBufferString(tc.reqBody),
			)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			resp := httptest.NewRecorder()

			server := gin.Default()
			udl := NewUserHandler(tc.getSvc(ctrl), nil)
			udl.RegisterRoutes(server)
			server.ServeHTTP(resp, req)

			assert.Equal(t, resp.Code, tc.wantCode)
			assert.Equal(t, resp.Body.String(), tc.wantBody)
		})
	}
}
