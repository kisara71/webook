package dao

import (
	"context"
	"database/sql"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB

		wantErr error
		user    UserEntity
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `user_entities` .*").WillReturnResult(sqlmock.NewResult(1, 1))
				return db
			},
			user: UserEntity{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			require.NoError(t, err)

			d := NewDao(db)
			err = d.Insert(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}
