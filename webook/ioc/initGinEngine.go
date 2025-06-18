package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/web/middleware"
	"github.com/kisara71/WeBook/webook/internal/web/oauth2"
	"github.com/kisara71/WeBook/webook/internal/web/user"
	"strings"
)

func InitGinEngine(mdw []gin.HandlerFunc, udl *user.Handler, oauth2Hdls []oauth2.Handler) *gin.Engine {
	server := gin.Default()
	server.Use(mdw...)
	udl.RegisterRoutes(server)
	for _, hd := range oauth2Hdls {
		hd.RegisterRoutes(server)
	}
	return server
}

func InitMiddleWare() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowMethods:     []string{"POST", "GET"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost") || strings.Contains(origin, "kisara71.xyz")
			},
		}),
		middleware.NewLoginJwtVerMiddleWare([]string{
			"/users/login_sms/code/send",
			"/users/login",
			"/users/login_sms",
			"/oauth2/wechat/authurl",
			"/oauth2/wechat/callback",
			"/oauth2/github/authurl",
			"/oauth2/github/callback",
		}).Build(),
	}
}
