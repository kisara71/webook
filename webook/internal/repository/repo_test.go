package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	cachemock "github.com/kisara71/WeBook/webook/internal/repository/cache/mocks"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
	daomock "github.com/kisara71/WeBook/webook/internal/repository/dao/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserRepositoryV1_FindById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.Dao, cache.UserCache)

		wantId  int64
		wantErr error
		Ctime   int64
		id      int64
	}{
		{
			name: "查找缓存成功",
			mock: func(ctrl *gomock.Controller) (dao.Dao, cache.UserCache) {
				d := daomock.NewMockDao(ctrl)
				uc := cachemock.NewMockUserCache(ctrl)

				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{Id: 1}, nil)
				return d, uc
			},

			wantId:  1,
			id:      1,
			wantErr: nil,
		},
		{
			name: "缓存未命中，但查找成功",
			mock: func(ctrl *gomock.Controller) (dao.Dao, cache.UserCache) {
				d := daomock.NewMockDao(ctrl)
				uc := cachemock.NewMockUserCache(ctrl)

				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotFound)
				uc.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)
				d.EXPECT().FindById(gomock.Any(), int64(1)).Return(domain.User{Id: 1}, nil)
				return d, uc
			},

			wantId:  1,
			id:      1,
			wantErr: nil,
		},
		{
			name: "查找失败",
			mock: func(ctrl *gomock.Controller) (dao.Dao, cache.UserCache) {
				d := daomock.NewMockDao(ctrl)
				uc := cachemock.NewMockUserCache(ctrl)

				uc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(domain.User{}, cache.ErrKeyNotFound)
				d.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("internal error"))
				return d, uc
			},

			wantErr: errors.New("internal error"),
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewUserRepository(ts.mock(ctrl))
			u, err := repo.FindById(context.Background(), ts.id)

			assert.Equal(t, u.Id, ts.wantId)
			assert.Equal(t, err, ts.wantErr)
		})
	}
}

func TestUserRepositoryV1_FindOrCreateByPhone(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) dao.Dao

		phone     string
		wantErr   error
		wantPhone string
	}{
		{
			name: "查找成功（非创建）",
			mock: func(ctrl *gomock.Controller) dao.Dao {
				d := daomock.NewMockDao(ctrl)
				d.EXPECT().FindUser(gomock.Any(), "Phone", "123").Return(domain.User{Phone: "123"}, nil)
				return d
			},
			phone:     "123",
			wantErr:   nil,
			wantPhone: "123",
		},
		{
			name: "查找成功（创建）",
			mock: func(ctrl *gomock.Controller) dao.Dao {
				d := daomock.NewMockDao(ctrl)
				d.EXPECT().FindUser(gomock.Any(), "Phone", "123").Return(domain.User{}, dao.ErrRecordNotFound)
				d.EXPECT().Insert(gomock.Any(), dao.UserEntity{
					Phone: sql.NullString{
						String: "123",
						Valid:  true,
					},
				}).Return(nil)
				d.EXPECT().FindUser(gomock.Any(), "Phone", "123").Return(domain.User{Phone: "123"}, nil)
				return d
			},
			phone:     "123",
			wantErr:   nil,
			wantPhone: "123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewUserRepository(tc.mock(ctrl), nil)

			u, err := repo.FindOrCreateByPhone(context.Background(), tc.phone)
			assert.Equal(t, u.Phone, tc.wantPhone)
			assert.Equal(t, err, tc.wantErr)

		})
	}
}
