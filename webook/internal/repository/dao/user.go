package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserPO struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique,uniqueIndex"`
	Password string

	Ctime int64
	Utime int64
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u UserPO) error {

	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	return dao.db.WithContext(ctx).Create(&u).Error
}
