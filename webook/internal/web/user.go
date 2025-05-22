package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/service"
	"net/http"
)

type UserHandler struct {
	regValidateEmail *regexp.Regexp
	regValidatePWD   *regexp.Regexp
	svc              *service.UserService
}

func InitUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		regValidateEmail: regexp.MustCompile("^[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*@[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*\\.[a-z]{2,}$", regexp.None),
		regValidatePWD:   regexp.MustCompile("^(?![0-9]+$)(?![a-zA-Z]+$)(?![0-9a-zA-Z]+$)(?![0-9\\W]+$)(?![a-zA-Z\\W]+$)[0-9A-Za-z\\W]{6,18}$", regexp.None),
		svc:              svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.signUp)
	ug.POST("/login", u.login)
	ug.POST("edit", u.edit)
	ug.GET("/profile", u.profile)

}

func (u *UserHandler) signUp(ctx *gin.Context) {
	type signUpReq struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		ConfirmPwd string `json:"confirmPassword"`
	}
	var req signUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Password != req.ConfirmPwd {
		ctx.String(http.StatusOK, "两次密码不同")
		return
	}
	//	validate email
	ok, err := u.regValidateEmail.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "错误的邮箱格式")
		return
	}
	//	validate password

	ok, err = u.regValidatePWD.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "错误的密码格式")
		return
	}
	//	service create user

	err = u.svc.Create(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v\n", req)
}

func (u *UserHandler) login(ctx *gin.Context) {
	type loginReq struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	req := loginReq{}
	if err := ctx.Bind(&req); err != nil {
		return
	}
}

func (u *UserHandler) profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "xxx")

}
func (u *UserHandler) edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "xxx")
}
