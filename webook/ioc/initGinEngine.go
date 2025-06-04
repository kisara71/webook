package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/web"
	"github.com/kisara71/WeBook/webook/internal/web/middleware"
	"strings"
)

func InitGinEngine(mdw []gin.HandlerFunc, udl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdw...)
	udl.RegisterRoutes(server)
	return server
}

func InitMiddleWare() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowMethods:     []string{"POST", "GET"},
			AllowHeaders:     []string{"Authorization", "Content-Type"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost") || strings.Contains(origin, "kisara71.xyz")
			},
		}),
		middleware.NewLoginJwtVerMiddleWare([]string{}).Build(),
	}
}
