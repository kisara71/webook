package user

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/GoTemplate/pkg/kstring"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"github.com/kisara71/WeBook/webook/internal/service/code_service"
	"github.com/kisara71/WeBook/webook/internal/service/user_service"
	"github.com/kisara71/WeBook/webook/internal/web"
	"github.com/kisara71/WeBook/webook/internal/web/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Handler struct {
	regValidateEmail *regexp.Regexp
	regValidatePWD   *regexp.Regexp
	regValidatePhone *regexp.Regexp
	userService      user_service.UserService
	smsService       code_service.CodeService
}

func NewUserHandler(userSvc user_service.UserService, smsSvc code_service.CodeService) *Handler {
	return &Handler{
		regValidateEmail: regexp.MustCompile("^[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*@[a-zA-Z0-9]+([-_.][a-zA-Z0-9]+)*\\.[a-z]{2,}$", regexp.None),
		regValidatePWD:   regexp.MustCompile("^(?![0-9]+$)(?![a-zA-Z]+$)(?![0-9a-zA-Z]+$)(?![0-9\\W]+$)(?![a-zA-Z\\W]+$)[0-9A-Za-z\\W]{6,18}$", regexp.None),
		regValidatePhone: regexp.MustCompile("/^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$/", regexp.None),
		userService:      userSvc,
		smsService:       smsSvc,
	}
}

func (u *Handler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.signUp)
	//ug.POST("/login", u.login)
	ug.POST("login", u.loginJwtVer)
	//ug.POST("edit", u.edit)
	//ug.GET("/profile", u.profile)
	ug.GET("/profile", u.profileJwtVer)
	ug.POST("/edit", u.editJwtVer)
	ug.POST("/login_sms/code/send", u.loginSmsSendCode)
	ug.POST("/login_sms", u.loginSms)

}

func (u *Handler) signUp(ctx *gin.Context) {
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

	err = u.userService.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, user_service.ErrEmailDuplicate) {
			ctx.String(http.StatusOK, "邮箱已注册")
			return
		}
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *Handler) login(ctx *gin.Context) {
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
	user, err = u.userService.FindUserByEmail(ctx, req.Email)
	if errors.Is(err, user_service.ErrInvalidEmailOrPassword) {
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

func (u *Handler) profile(ctx *gin.Context) {
	userId := sessions.Default(ctx).Get("userId").(int64)
	var (
		err  error
		user domain.User
	)
	user, err = u.userService.FindUserById(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "无效的帐号")
		return
	}
	//user, _ = u.userService.FindUserInfoById(ctx, userId)

	ctx.JSON(http.StatusOK, gin.H{
		"Email":    user.Email,
		"Phone":    "",
		"AboutMe":  user.AboutMe,
		"Nickname": user.Name,
		"Birthday": user.Birthday,
	})
}
func (u *Handler) edit(ctx *gin.Context) {
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

	err = u.userService.Edit(ctx, domain.User{
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

func (u *Handler) loginJwtVer(ctx *gin.Context) {
	type loginJwtReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var loginReq loginJwtReq

	if err := ctx.Bind(&loginReq); err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	user, err := u.userService.Login(ctx, loginReq.Email, loginReq.Password)
	if err != nil {
		if errors.Is(err, user_service.ErrInvalidEmailOrPassword) {
			ctx.String(http.StatusOK, "密码或邮箱错误")
			return
		}
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	err = util.SetJwtToken(ctx, user.Id)
	if err != nil {
		// log
	}
	ctx.String(http.StatusOK, "login successfully")
}

func (u *Handler) profileJwtVer(ctx *gin.Context) {
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
	user, err = u.userService.FindUserById(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "无效的帐号")
		return
	}
	//user, _ = u.userService.FindUserInfoById(ctx, userId)

	ctx.JSON(http.StatusOK, gin.H{
		"Email":    user.Email,
		"Phone":    user.Phone,
		"AboutMe":  user.AboutMe,
		"Nickname": user.Name,
		"Birthday": user.Birthday,
	})
}

func (u *Handler) editJwtVer(ctx *gin.Context) {
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

	err = u.userService.Edit(ctx, domain.User{
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
func (u *Handler) loginSmsSendCode(ctx *gin.Context) {
	type SendCodeReq struct {
		Phone string `json:"phone"`
	}
	var req SendCodeReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	err = u.smsService.Send(ctx, "user", req.Phone)
	switch {
	case errors.Is(err, nil):
		ctx.JSON(http.StatusOK, web.Result{
			Code: 0,
			Msg:  "发送成功",
		})
	case errors.Is(err, code_service.ErrSendTooFrequent):
		ctx.JSON(http.StatusOK, web.Result{
			Code: 4,
			Msg:  "验证码发送太频繁",
		})
	default:
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *Handler) loginSms(ctx *gin.Context) {
	type PhoneLoginReq struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req PhoneLoginReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}
	code, err := kstring.ToInt(req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}
	same, err := u.smsService.VerifyCode(ctx, "user", req.Phone, code)
	if same {
		newUD, err := u.userService.FindOrCreateByPhone(ctx, req.Phone)
		if err != nil {
			ctx.JSON(http.StatusOK, web.Result{
				Code: 4,
				Msg:  "system error",
			})
			return
		}
		err = util.SetJwtToken(ctx, newUD.Id)
		if err != nil {
			//	log
		}
		ctx.JSON(http.StatusOK, web.Result{
			Msg: "登录成功",
		})
		return
	} else {
		switch {
		case errors.Is(err, code_service.ErrWrongCode):
			ctx.JSON(http.StatusOK, web.Result{
				Code: 4,
				Msg:  "验证码错误，请重新尝试",
			})
		case errors.Is(err, code_service.ErrTooManyVerifications):
			ctx.JSON(http.StatusOK, web.Result{
				Code: 4,
				Msg:  "错误次数过多，请重新获取验证码",
			})
		default:
			ctx.JSON(http.StatusOK, web.Result{
				Code: 4,
				Msg:  "系统错误",
			})
		}

	}
}
