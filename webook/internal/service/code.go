package service

import (
	"context"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/repository"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"math/rand/v2"
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error)
}

func NewCodeService(cp repository.CodeRepository, svc sms.Service) CodeService {
	return newCodeServiceV1(cp, svc)
}

var (
	ErrSendTooFrequent      = repository.ErrSendTooFrequent
	ErrSystemError          = repository.ErrSystemError
	ErrInvalidCode          = repository.ErrInvalidCode
	ErrTooManyVerifications = repository.ErrTooManyVerifications
	ErrWrongCode            = repository.ErrWrongCode
)

type codeServiceV1 struct {
	codeRepo repository.CodeRepository
	sms      sms.Service
}

func newCodeServiceV1(r repository.CodeRepository, sms sms.Service) CodeService {
	return &codeServiceV1{
		codeRepo: r,
		sms:      sms,
	}
}

func (C *codeServiceV1) Send(ctx context.Context, biz, phone string) error {
	code := C.generateCode()
	err := C.codeRepo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = C.sms.Send(ctx, sms.Message{
		TemplateParm: fmt.Sprintf("{\"code\":\"%d\"}", code),
		PhoneNumbers: phone,
	})
	return err
}

func (C *codeServiceV1) VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error) {
	return C.codeRepo.VerifyCode(ctx, biz, phone, code)
}

func (C *codeServiceV1) generateCode() int {
	return rand.IntN(100000) + 100000
}
