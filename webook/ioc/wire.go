//go:build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kisara71/WeBook/webook/internal/repository"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
	"github.com/kisara71/WeBook/webook/internal/service"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"github.com/kisara71/WeBook/webook/internal/web"
)

func InitWebServer() *gin.Engine {
	wire.Build(

		InitDatabase, InitRedis,

		dao.NewDao,

		cache.NewUserCache, cache.NewCodeCache,

		repository.NewCodeRepository, repository.NewUserRepository,

		sms.NewSMSService,

		service.NewCodeService, service.NewUserService,

		web.NewUserHandler,

		InitMiddleWare,

		InitGinEngine,
	)
	return new(gin.Engine)
}
