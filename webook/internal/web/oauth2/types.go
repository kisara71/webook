package oauth2

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutes(server *gin.Engine)
	AuthURL(ctx *gin.Context)
	CallBack(ctx *gin.Context)
}
