package ioc

import (
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"github.com/kisara71/WeBook/webook/internal/service/sms/decorator"
	"github.com/kisara71/WeBook/webook/internal/service/sms/provider"
	"github.com/kisara71/WeBook/webook/pkg/ratelimit"
)

func initSMS(limiter ratelimit.RateLimiter) sms.Service {
	basicSMS := provider.NewSMSMemory()
	limitSMS := decorator.NewRateLimitSMSService(basicSMS, limiter)
	failoverSMS := decorator.NewFailOverSMSService([]sms.Service{limitSMS})
	return failoverSMS
}
