package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"gorm.io/gorm"
	"time"
)

type Dao interface {
	Insert(ctx context.Context, u UserEntity) error
	Edit(ctx context.Context, info domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindUser(ctx context.Context, filed string, value any) (domain.User, error)
}

func NewDao(db *gorm.DB) Dao {
	return newGormDao(db)
}

type UserEntity struct {
	Id       int64          `gorm:"primaryKey;autoIncrement"`
	Email    sql.NullString `gorm:"uniqueIndex;type:varchar(50)"`
	Password string         `gorm:"type:varchar(200)"`
	Phone    sql.NullString `gorm:"uniqueIndex;type:varchar(20)"`
	Name     string         `gorm:"type:varchar(10)"`
	Birthday string         `gorm:"type:date;default:NULL"`
	AboutMe  string         `gorm:"varchar(50)"`
	Ctime    int64
	Utime    int64
}

const (
	uniqueConflictsErrno uint16 = 1062
)

var (
	ErrEmailDuplicate         = errors.New("邮箱冲突")
	ErrInvalidEmailOrPassword = gorm.ErrRecordNotFound
	ErrRecordNotFound         = gorm.ErrRecordNotFound
)

type gormDao struct {
	db *gorm.DB
}

func newGormDao(db *gorm.DB) Dao {
	return &gormDao{
		db: db,
	}
}

func (gd *gormDao) Insert(ctx context.Context, u UserEntity) error {

	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := gd.db.WithContext(ctx).Create(&u).Error

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == uniqueConflictsErrno {
			return ErrEmailDuplicate
		} else {
			return err
		}
	}
	return nil
}

func (gd *gormDao) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var u UserEntity
	err := gd.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, ErrInvalidEmailOrPassword) {
			return domain.User{}, ErrInvalidEmailOrPassword
		} else {
			return domain.User{}, err
		}
	}
	return gd.entityToDomain(u), nil
}

func (gd *gormDao) Edit(ctx context.Context, info domain.User) error {
	return gd.db.WithContext(ctx).Where("Id = ?", info.Id).Updates(&UserEntity{
		Name:     info.Name,
		Birthday: info.Birthday,
		AboutMe:  info.AboutMe,
		Utime:    time.Now().UnixMilli(),
	}).Error
}

//	func (dao *gormDao) FindUserInfoById(ctx context.Context, id int64) (domain.User, error) {
//		var u dao.UserEntity
//		err := dao.db.WithContext(ctx).Where("Id = ?", id).First(&u).Error
//		return domain.User{
//			Id:       u.Id,
//			Birthday: u.Birthday,
//			Name:     u.Name,
//			AboutMe:  u.AboutMe,
//		}, err
//	}
func (gd *gormDao) FindById(ctx context.Context, id int64) (domain.User, error) {
	var u UserEntity
	err := gd.db.WithContext(ctx).Where("Id = ?", id).First(&u).Error
	return gd.entityToDomain(u), err
}
func (gd *gormDao) FindUser(ctx context.Context, filed string, value any) (domain.User, error) {
	var u UserEntity
	err := gd.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", filed), value).First(&u).Error
	return gd.entityToDomain(u), err
}

func (gd *gormDao) entityToDomain(ud UserEntity) domain.User {
	return domain.User{
		Id:       ud.Id,
		Name:     ud.Name,
		Phone:    ud.Phone.String,
		Email:    ud.Email.String,
		AboutMe:  ud.AboutMe,
		Birthday: ud.Birthday,
	}
}
