package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/service/auth_binding_service"
	"github.com/kisara71/WeBook/webook/internal/service/oauth2"
	"github.com/kisara71/WeBook/webook/internal/web"
	"github.com/kisara71/WeBook/webook/internal/web/util"
	"net/http"
)

type GithubHandler struct {
	githubSvc      oauth2.Service
	authBindingSvc auth_binding_service.Service
}

func NewGithubHandler(oauth2Svc oauth2.Service, authBindingSvc auth_binding_service.Service) Handler {
	return &GithubHandler{
		githubSvc:      oauth2Svc,
		authBindingSvc: authBindingSvc,
	}
}

func (g *GithubHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("oauth2/github/authurl", g.AuthURL)
	server.GET("oauth2/github/callback", g.CallBack)
}

func (g *GithubHandler) AuthURL(ctx *gin.Context) {
	url := g.githubSvc.AuthURL(ctx)

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Data: url,
	})
}

func (g *GithubHandler) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	info, err := g.githubSvc.ExchangeCode(ctx, code, state)

	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	info, err = g.authBindingSvc.FindOrCreateOauth2Binding(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = util.SetJwtToken(ctx, info.UserID)
	if err != nil {
		// log
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "登录成功",
	})
}
