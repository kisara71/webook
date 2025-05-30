package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"gorm.io/gorm"
	"time"
)

const (
	uniqueConflictsErrno uint16 = 1062
)

var (
	ErrEmailDuplicate         = errors.New("邮箱冲突")
	ErrInvalidEmailOrPassword = gorm.ErrRecordNotFound
	ErrIdNotFound             = gorm.ErrRecordNotFound
)

type UserPO struct {
	Id       int64  `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"uniqueIndex;type:varchar(50)"`
	Password string `gorm:"type:varchar(200)"`

	Name     string `gorm:"type:varchar(10)"`
	Birthday string `gorm:"type:date;default:NULL"`
	AboutMe  string `gorm:"varchar(50)"`
	Ctime    int64
	Utime    int64
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

	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == uniqueConflictsErrno {
			return ErrEmailDuplicate
		}
	}
	return nil
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var u UserPO
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u)
	if err.Error != nil {
		if errors.Is(err.Error, ErrInvalidEmailOrPassword) {
			return domain.User{}, ErrInvalidEmailOrPassword
		} else {
			return domain.User{}, err.Error
		}
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (dao *UserDao) Edit(ctx context.Context, info domain.User) error {
	return dao.db.WithContext(ctx).Where("Id = ?", info.Id).Updates(&UserPO{
		Name:     info.Name,
		Birthday: info.Birthday,
		AboutMe:  info.AboutMe,
		Utime:    time.Now().UnixMilli(),
	}).Error
}

//	func (dao *UserDao) FindUserInfoById(ctx context.Context, id int64) (domain.User, error) {
//		var u UserPO
//		err := dao.db.WithContext(ctx).Where("Id = ?", id).First(&u).Error
//		return domain.User{
//			Id:       u.Id,
//			Birthday: u.Birthday,
//			Name:     u.Name,
//			AboutMe:  u.AboutMe,
//		}, err
//	}
func (dao *UserDao) FindUserById(ctx context.Context, id int64) (domain.User, error) {
	var u UserPO
	err := dao.db.WithContext(ctx).Where("Id = ?", id).First(&u).Error
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Birthday: u.Birthday,
		Name:     u.Name,
		AboutMe:  u.AboutMe,
	}, err
}
