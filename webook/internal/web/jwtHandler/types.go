package jwtHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetJwtToken(ctx *gin.Context, id int64) error
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64
}
