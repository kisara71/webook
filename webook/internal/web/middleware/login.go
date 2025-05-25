package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddleBuilder struct {
	ignorePath []string
}

func NewLoginMiddlerBuilder() *LoginMiddleBuilder {
	return &LoginMiddleBuilder{
		ignorePath: []string{
			"/users/signup",
			"/users/login",
		},
	}
}

func (l *LoginMiddleBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.ignorePath {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		if id := sess.Get("userId"); id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
