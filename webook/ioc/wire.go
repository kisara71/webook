//go:build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kisara71/WeBook/webook/internal/repository"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
	"github.com/kisara71/WeBook/webook/internal/service"
	"github.com/kisara71/WeBook/webook/internal/web"
	"github.com/kisara71/WeBook/webook/internal/web/jwtHandler"
)

func InitWebServer() *gin.Engine {
	wire.Build(

		InitDatabase, InitRedis, jwtHandler.NewJwtHandler, initRateLimiter,

		dao.NewDao,

		cache.NewUserCache, cache.NewCodeCache,

		repository.NewCodeRepository, repository.NewUserRepository,

		initSMS,

		service.NewCodeService, service.NewUserService, initWeChatService,

		web.NewUserHandler, web.NewWeChatHandler,

		InitMiddleWare,

		InitGinEngine,
	)
	return new(gin.Engine)
}
