package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/service/auth_binding_service"
	"github.com/kisara71/WeBook/webook/internal/service/oauth2"
	"github.com/kisara71/WeBook/webook/internal/web"
	"github.com/kisara71/WeBook/webook/internal/web/util"
	"net/http"
)

type WeChatHandler struct {
	wechatSvc      oauth2.Service
	authBindingSvc auth_binding_service.Service
}

func NewWeChatHandler(oauth2Svc oauth2.Service, authBindingSvc auth_binding_service.Service) Handler {
	return &WeChatHandler{
		wechatSvc:      oauth2Svc,
		authBindingSvc: authBindingSvc,
	}
}

func (w *WeChatHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("oauth2/wechat/authurl", w.AuthURL)
	server.GET("oauth2/wechat/callback", w.CallBack)
}
func (w *WeChatHandler) AuthURL(ctx *gin.Context) {
	url := w.wechatSvc.AuthURL(ctx)

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Data: url,
	})

}

func (w *WeChatHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	info, err := w.wechatSvc.ExchangeCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	info, err = w.authBindingSvc.FindOrCreateOauth2Binding(ctx, info)
	if err != nil {
		return
	}
	err = util.SetJwtToken(ctx, info.UserID)
}
