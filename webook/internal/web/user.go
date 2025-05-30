package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
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
	//ug.POST("/login", u.login)
	ug.POST("login", u.loginJwtVer)
	//ug.POST("edit", u.edit)
	//ug.GET("/profile", u.profile)
	ug.GET("/profile", u.profileJwtVer)
	ug.POST("edit", u.editJwtVer)

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
		if errors.Is(err, service.ErrEmailDuplicate) {
			ctx.String(http.StatusOK, "邮箱已注册")
			return
		}
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := loginReq{}
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	var user domain.User
	user, err = u.svc.FindByEmail(ctx, req.Email)
	if errors.Is(err, service.ErrInvalidEmailOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}

	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	err = session.Save()
	ctx.String(http.StatusOK, "登录成功")

}

func (u *UserHandler) profile(ctx *gin.Context) {
	userId := sessions.Default(ctx).Get("userId").(int64)
	var (
		err  error
		user domain.User
	)
	user, err = u.svc.FindUserById(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "无效的帐号")
		return
	}
	//user, _ = u.svc.FindUserInfoById(ctx, userId)

	ctx.JSON(http.StatusOK, gin.H{
		"Email":    user.Email,
		"Phone":    "",
		"AboutMe":  user.AboutMe,
		"Nickname": user.Name,
		"Birthday": user.Birthday,
	})
}
func (u *UserHandler) edit(ctx *gin.Context) {
	type editReq struct {
		NickName string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	userId := sessions.Default(ctx).Get("userId").(int64)

	var req editReq
	var err error
	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	err = u.svc.Edit(ctx, domain.User{
		Id:       userId,
		Name:     req.NickName,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "修改成功",
	})
}

func (u *UserHandler) loginJwtVer(ctx *gin.Context) {
	type loginJwtReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var loginReq loginJwtReq

	if err := ctx.Bind(&loginReq); err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	du, err := u.svc.FindByEmail(ctx, loginReq.Email)
	if err != nil {
		ctx.String(http.StatusOK, "invalid email")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(du.Password), []byte(loginReq.Password))
	if err != nil {
		ctx.String(http.StatusOK, "check password or email")
		return
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		UserId: du.Id,
	}).SignedString([]byte("2yJPXiYFxjQC6D4G73vHKoJ90bv7DNixOIsTDdulApdjv0QNoK5rOL9xSASLlQvg"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal error")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "login successfully")
}

func (u *UserHandler) profileJwtVer(ctx *gin.Context) {
	ctxMsg, ok := ctx.Get("userId")
	userId, _ := ctxMsg.(int64)
	if !ok {
		ctx.String(http.StatusUnauthorized, "invalid login token")
		return
	}
	var (
		err  error
		user domain.User
	)
	user, err = u.svc.FindUserById(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "无效的帐号")
		return
	}
	//user, _ = u.svc.FindUserInfoById(ctx, userId)

	ctx.JSON(http.StatusOK, gin.H{
		"Email":    user.Email,
		"Phone":    "",
		"AboutMe":  user.AboutMe,
		"Nickname": user.Name,
		"Birthday": user.Birthday,
	})
}

func (u *UserHandler) editJwtVer(ctx *gin.Context) {
	type editReq struct {
		NickName string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	ctxMsg, ok := ctx.Get("userId")
	if !ok {
		ctx.String(http.StatusUnauthorized, "invalid token")
		return
	}
	userId, _ := ctxMsg.(int64)
	var req editReq
	var err error
	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	err = u.svc.Edit(ctx, domain.User{
		Id:       userId,
		Name:     req.NickName,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "修改成功",
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64
}
