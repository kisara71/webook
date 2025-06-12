package cache

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/kisara71/WeBook/webook/internal/domain"
	cachemock "github.com/kisara71/WeBook/webook/internal/repository/cache/mocks"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestRedisUserCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable
		in   domain.User

		wantErr error
	}{
		{
			name: "缓存设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				r := cachemock.NewMockCmdable(ctrl)
				r.EXPECT().Set(context.Background(), gomock.Any(), gomock.Any(), time.Minute*15).Return(redis.NewStatusCmd(context.Background()))
				return r
			},
			in: domain.User{
				Id: 1,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewUserCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.in)

			assert.Equal(t, err, tc.wantErr)
		})
	}
}
