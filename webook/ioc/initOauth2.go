package ioc

import (
	oauth2Config "github.com/kisara71/WeBook/webook/config/oauth2"
	"github.com/kisara71/WeBook/webook/internal/service/auth_binding_service"
	oauthSvc "github.com/kisara71/WeBook/webook/internal/service/oauth2"
	"github.com/kisara71/WeBook/webook/internal/web/oauth2"
)

func initOauth2Handlers(authBindingSvc auth_binding_service.Service) []oauth2.Handler {
	handlers := make([]oauth2.Handler, 0)

	handlers = append(handlers, oauth2.NewWeChatHandler(initWechatService(), authBindingSvc))
	handlers = append(handlers, oauth2.NewGithubHandler(initGithubService(), authBindingSvc))
	return handlers
}

func initWechatService() oauthSvc.Service {
	return oauthSvc.NewWechatService(oauth2Config.Configs["Wechat"])
}

func initGithubService() oauthSvc.Service {
	return oauthSvc.NewGithubService(oauth2Config.Configs["Github"])
}
