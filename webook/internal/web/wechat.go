package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/service"
	"github.com/kisara71/WeBook/webook/internal/web/jwtHandler"
	"net/http"
)

type WeChatHandler struct {
	wechatSvc service.WechatService
	userSvc   service.UserService
	jwtHd     jwtHandler.Handler
}

func NewWeChatHandler(wechatSvc service.WechatService, userSvc service.UserService, jwtHd jwtHandler.Handler) *WeChatHandler {
	return &WeChatHandler{
		wechatSvc: wechatSvc,
		userSvc:   userSvc,
		jwtHd:     jwtHd,
	}
}

func (w *WeChatHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("oauth2/wechat/authurl", w.AuthURL)
	server.GET("oauth2/wechat/callback", w.CallBack)
}
func (w *WeChatHandler) AuthURL(ctx *gin.Context) {
	url, err := w.wechatSvc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "AuthURL构造失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Data: url,
	})

}

func (w *WeChatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	info, err := w.wechatSvc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	user, err := w.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		return
	}
	err = w.jwtHd.SetJwtToken(ctx, user.Id)
}
