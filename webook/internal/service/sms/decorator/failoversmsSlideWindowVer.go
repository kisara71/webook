package decorator

import (
	"context"
	"errors"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"github.com/kisara71/WeBook/webook/internal/service/sms/decorator/util"
	"sync/atomic"
)

type FailOverSMSSlideWindowVer struct {
	svcs []sms.Service
	idx  int32
	sd   *util.FailOverSlideWindow
}

func NewFailOverSMSSlideWindowVer(sd *util.FailOverSlideWindow, svcs []sms.Service) *FailOverSMSSlideWindowVer {
	return &FailOverSMSSlideWindowVer{
		svcs: svcs,
		sd:   sd,
		idx:  0,
	}
}

func (f *FailOverSMSSlideWindowVer) Send(ctx context.Context, msg sms.Message) error {

	idx := atomic.LoadInt32(&f.idx)

	svc := f.svcs[idx]
	err := svc.Send(ctx, msg)
	f.sd.Add(err == nil)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrLimited):
		// async add msg

	case f.sd.ShouldFailOver():
		newIdx := (idx + 1) % int32(len(f.svcs))
		atomic.CompareAndSwapInt32(&f.idx, idx, newIdx)

	}
	return err
}
