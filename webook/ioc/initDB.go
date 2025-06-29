package ioc

import (
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&dao.UserEntity{}, &dao.Oauth2BindingEntity{})
	if err != nil {
		panic(err)
	}
	return db
}
