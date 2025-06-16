package code_service

import (
	"context"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/repository/code_repo"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"math/rand/v2"
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error)
}

func NewCodeService(cp code_repo.CodeRepository, svc sms.Service) CodeService {
	return newCodeServiceV1(cp, svc)
}

var (
	ErrSendTooFrequent      = code_repo.ErrSendTooFrequent
	ErrSystemError          = code_repo.ErrSystemError
	ErrInvalidCode          = code_repo.ErrInvalidCode
	ErrTooManyVerifications = code_repo.ErrTooManyVerifications
	ErrWrongCode            = code_repo.ErrWrongCode
)

type codeServiceV1 struct {
	codeRepo code_repo.CodeRepository
	sms      sms.Service
}

func newCodeServiceV1(r code_repo.CodeRepository, sms sms.Service) CodeService {
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
