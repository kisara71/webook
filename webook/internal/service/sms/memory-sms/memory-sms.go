package memory_sms

import (
	"context"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Send(ctx context.Context, message sms.Message) error {
	fmt.Println(message)
	return nil
}
