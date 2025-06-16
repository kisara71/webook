package decorator

import (
	"context"
	"errors"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"github.com/kisara71/WeBook/webook/pkg/ratelimit"
)

var (
	ErrLimited = errors.New("limited")
)

type RateLimitSMSService struct {
	sms.Service
	limiter ratelimit.RateLimiter
}

func NewRateLimitSMSService(svc sms.Service, limiter ratelimit.RateLimiter) sms.Service {
	return &RateLimitSMSService{
		limiter: limiter,
		Service: svc,
	}
}

func (r *RateLimitSMSService) Send(ctx context.Context, msg sms.Message) error {
	limited, err := r.limiter.Limit(ctx, "sms_service")
	if err != nil {
		return fmt.Errorf("sms 限流检测失败%w", err)

	}
	if limited {
		return ErrLimited
	}
	return r.Service.Send(ctx, msg)
}
