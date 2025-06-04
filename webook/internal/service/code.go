package service

import (
	"context"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/repository"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"math/rand/v2"
)

var (
	ErrSendTooFrequent      = repository.ErrSendTooFrequent
	ErrSystemError          = repository.ErrSystemError
	ErrInvalidCode          = repository.ErrInvalidCode
	ErrTooManyVerifications = repository.ErrTooManyVerifications
	ErrWrongCode            = repository.ErrWrongCode
)

type CodeService struct {
	codeRepo *repository.CodeRepository
	sms      sms.Service
}

func NewCodeService(r *repository.CodeRepository, sms sms.Service) *CodeService {
	return &CodeService{
		codeRepo: r,
		sms:      sms,
	}
}

func (c *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := c.generateCode()
	err := c.codeRepo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = c.sms.Send(ctx, sms.Message{
		TemplateParm: fmt.Sprintf("{\"code\":\"%d\"}", code),
		PhoneNumbers: phone,
	})
	return err
}

func (c *CodeService) VerifyCode(ctx context.Context, biz, phone string, code int) (bool, error) {
	return c.codeRepo.VerifyCode(ctx, biz, phone, code)
}

func (c *CodeService) generateCode() int {
	return rand.IntN(100000) + 100000
}
