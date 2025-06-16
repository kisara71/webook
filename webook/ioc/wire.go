//go:build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kisara71/WeBook/webook/internal/repository/auth_binding_repo"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/code_repo"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
	"github.com/kisara71/WeBook/webook/internal/repository/user_repo"
	"github.com/kisara71/WeBook/webook/internal/service/auth_binding_service"
	"github.com/kisara71/WeBook/webook/internal/service/code_service"
	"github.com/kisara71/WeBook/webook/internal/service/user_service"
	"github.com/kisara71/WeBook/webook/internal/web/user"
)

func InitWebServer() *gin.Engine {
	wire.Build(

		InitDatabase, InitRedis, initRateLimiter,

		dao.NewDao,

		cache.NewUserCache, cache.NewCodeCache,

		code_repo.NewCodeRepository, user_repo.NewUserRepository, auth_binding_repo.NewRepository,

		initSMS, auth_binding_service.NewService,

		code_service.NewCodeService, user_service.NewUserService,

		user.NewUserHandler, initOauth2Handlers,

		InitMiddleWare,

		InitGinEngine,
	)
	return new(gin.Engine)
}
