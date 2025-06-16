package ioc

import (
	"github.com/kisara71/WeBook/webook/internal/service/auth_binding_service"
	oauthSvc "github.com/kisara71/WeBook/webook/internal/service/oauth2"
	"github.com/kisara71/WeBook/webook/internal/web/oauth2"
)

func initOauth2Handlers(authBinding_Svc auth_binding_service.Service) []oauth2.Handler {
	hdls := make([]oauth2.Handler, 0)

	hdls = append(hdls, oauth2.NewWeChatHandler(initWechatService(oauthSvc.Config{
		ClientID:     "dfsdf",
		ClientSecret: "sfd",
	}), authBinding_Svc))
	return hdls
}

func initWechatService(config oauthSvc.Config) oauthSvc.Service {
	return oauthSvc.NewWechatService(config)
}
