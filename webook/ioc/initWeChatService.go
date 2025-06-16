package ioc

import "github.com/kisara71/WeBook/webook/internal/service"

func initWeChatService() service.WechatService {
	return service.NewWechatService("1233")
}
