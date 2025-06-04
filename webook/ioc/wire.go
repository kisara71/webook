//go:build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kisara71/WeBook/webook/internal/repository"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
	"github.com/kisara71/WeBook/webook/internal/service"
	"github.com/kisara71/WeBook/webook/internal/web"
)

func InitWebServer() *gin.Engine {
	wire.Build(

		InitDatabase, InitRedis,

		dao.NewUserDao,

		InitUserCache, InitCodeCache,

		repository.NewUserRepository, repository.NewCodeRepository,

		InitSMSService,

		service.NewUserService, service.NewCodeService,

		web.NewUserHandler,

		InitMiddleWare,

		InitGinEngine,
	)
	return new(gin.Engine)
}
