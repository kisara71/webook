package ioc

import (
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	memory_sms "github.com/kisara71/WeBook/webook/internal/service/sms/memory-sms"
)

func InitSMSService() sms.Service {
	return memory_sms.NewService()
}
