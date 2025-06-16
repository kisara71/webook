package provider

import (
	"context"
	"fmt"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
)

type MemoryService struct {
}

func NewSMSMemory() sms.Service {
	return &MemoryService{}
}
func (M *MemoryService) Send(ctx context.Context, message sms.Message) error {
	fmt.Println(message)
	return nil
}
