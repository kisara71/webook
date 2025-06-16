package util

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func SetJwtToken(ctx *gin.Context, id int64) error {
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		UserId: id,
	}).SignedString([]byte("2yJPXiYFxjQC6D4G73vHKoJ90bv7DNixOIsTDdulApdjv0QNoK5rOL9xSASLlQvg"))
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
