package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"gorm.io/gorm"
	"time"
)

type Dao interface {
	InsertUser(ctx context.Context, user domain.User) (domain.User, error)
	Edit(ctx context.Context, info domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindUser(ctx context.Context, filed string, value any) (domain.User, error)
	FindBinding(ctx context.Context, provider string, externalID string) (domain.Oauth2Binding, error)
	InsertOauth2Binding(ctx context.Context, binding domain.Oauth2Binding) (domain.Oauth2Binding, error)
}

func NewDao(db *gorm.DB) Dao {
	return newGormDao(db)
}

const (
	uniqueConflictsErrno uint16 = 1062
)

var (
	ErrDuplicate              = errors.New("用户已存在冲突")
	ErrInvalidEmailOrPassword = gorm.ErrRecordNotFound
	ErrRecordNotFound         = gorm.ErrRecordNotFound
	ErrSystemError            = errors.New("系统错误")
)

type gormDao struct {
	db *gorm.DB
}

func (gd *gormDao) InsertOauth2Binding(ctx context.Context, binding domain.Oauth2Binding) (domain.Oauth2Binding, error) {
	entity := oauth2BindingDomainToEntity(&binding)

	now := time.Now().UnixMilli()
	entity.Utime = now
	entity.Ctime = now
	err := gd.db.WithContext(ctx).Create(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.Oauth2Binding{}, ErrDuplicate
		} else {
			return domain.Oauth2Binding{}, ErrSystemError
		}
	}
	binding.UserID = entity.UserID
	binding.ID = entity.ID
	return binding, nil
}

func (gd *gormDao) FindBinding(ctx context.Context, provider string, externalID string) (domain.Oauth2Binding, error) {
	var res domain.Oauth2Binding
	err := gd.db.WithContext(ctx).Where("provider = ? AND externalID = ?", provider, externalID).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Oauth2Binding{}, ErrRecordNotFound
		} else {
			return domain.Oauth2Binding{}, ErrSystemError
		}
	}
	return res, nil
}

func newGormDao(db *gorm.DB) Dao {
	return &gormDao{
		db: db,
	}
}

func (gd *gormDao) InsertUser(ctx context.Context, user domain.User) (domain.User, error) {

	u := userDomainTOentity(&user)
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := gd.db.WithContext(ctx).Create(&u).Error

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == uniqueConflictsErrno {
			return domain.User{}, ErrDuplicate
		} else {
			// log
			return domain.User{}, ErrSystemError
		}
	}
	user.Id = u.ID
	return user, nil
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
	return userEntityToDomain(&u), nil
}

func (gd *gormDao) Edit(ctx context.Context, info domain.User) error {
	return gd.db.WithContext(ctx).Where("ID = ?", info.Id).Updates(&UserEntity{
		Name:     info.Name,
		Birthday: info.Birthday,
		AboutMe:  info.AboutMe,
		Utime:    time.Now().UnixMilli(),
	}).Error
}

//	func (dao *gormDao) FindUserInfoById(ctx context.Context, id int64) (domain.User, error) {
//		var u dao.UserEntity
//		err := dao.db.WithContext(ctx).Where("ID = ?", id).First(&u).Error
//		return domain.User{
//			ID:       u.ID,
//			Birthday: u.Birthday,
//			Name:     u.Name,
//			AboutMe:  u.AboutMe,
//		}, err
//	}
func (gd *gormDao) FindById(ctx context.Context, id int64) (domain.User, error) {
	var u UserEntity
	err := gd.db.WithContext(ctx).Where("ID = ?", id).First(&u).Error
	return userEntityToDomain(&u), err
}
func (gd *gormDao) FindUser(ctx context.Context, filed string, value any) (domain.User, error) {
	var u UserEntity
	err := gd.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", filed), value).First(&u).Error
	return userEntityToDomain(&u), err
}
