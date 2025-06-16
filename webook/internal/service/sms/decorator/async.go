package decorator

import (
	"context"
	"encoding/json"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"github.com/redis/go-redis/v9"
	"sync/atomic"
)

type AsyncMessage struct {
	PhoneNumbers string
	TemplateParm string
	Retry        int
}

type AsyncSMSService struct {
	ConsecutiveFailure int32
	maxRetry           int
	sms.Service
	cmd redis.Cmdable
	key string
}

func NewAsyncSMSService(svc sms.Service, cmd redis.Cmdable) sms.Service {
	return &AsyncSMSService{
		Service:            svc,
		cmd:                cmd,
		key:                "sms:async:queue",
		ConsecutiveFailure: 0,
		maxRetry:           3,
	}
}
func (a *AsyncSMSService) worker(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		for {
			msg, err := a.cmd.BLPop(ctx, 0, a.key).Result()
			if err != nil {
				// log
				continue
			}
			var val AsyncMessage
			err = json.Unmarshal([]byte(msg[1]), &val)
			if err != nil {
				// log
				continue
			}
			if val.Retry >= a.maxRetry {
				atomic.AddInt32(&a.ConsecutiveFailure, 1)
				continue
			}
			err = a.Service.Send(context.Background(), sms.Message{
				TemplateParm: val.TemplateParm,
				PhoneNumbers: val.PhoneNumbers,
			})
			if err != nil {
				//  log
				_ = a.Add(context.Background(), val)
			}
		}
	}
}

func (a *AsyncSMSService) Add(ctx context.Context, asyncMsg AsyncMessage) error {

	req, err := json.Marshal(asyncMsg)
	if err != nil {
		//  log
		return err
	}
	_, err = a.cmd.LPush(ctx, a.key, req).Result()
	return err
}
