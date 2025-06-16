package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kisara71/WeBook/webook/internal/web/util"
	"net/http"
	"strings"
	"time"
)

type LoginJwtVerMiddleWare struct {
	IgnorePath []string
}

func (l *LoginJwtVerMiddleWare) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.IgnorePath {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		au := ctx.GetHeader("Authorization")
		if au == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.Split(au, " ")
		var userclaim util.UserClaims
		tokenStr, err := jwt.ParseWithClaims(token[1], &userclaim, func(token *jwt.Token) (interface{}, error) {
			return []byte("2yJPXiYFxjQC6D4G73vHKoJ90bv7DNixOIsTDdulApdjv0QNoK5rOL9xSASLlQvg"), nil
		})
		if err != nil {
			ctx.String(http.StatusInternalServerError, "internal error")
			return
		}
		if !tokenStr.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if userclaim.ExpiresAt.Sub(time.Now()) < time.Minute*1 {
			newToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS512, &util.UserClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
				},
				UserId: userclaim.UserId,
			}).SignedString([]byte("2yJPXiYFxjQC6D4G73vHKoJ90bv7DNixOIsTDdulApdjv0QNoK5rOL9xSASLlQvg"))
			ctx.Header("x-jwtHandler-token", newToken)
		}
		ctx.Set("userId", userclaim.UserId)
	}
}

func NewLoginJwtVerMiddleWare(ignorePath []string) *LoginJwtVerMiddleWare {
	if len(ignorePath) == 0 {
		return nil
	}
	return &LoginJwtVerMiddleWare{
		IgnorePath: ignorePath,
	}
}
