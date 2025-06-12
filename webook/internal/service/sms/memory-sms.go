package sms

import (
	"context"
	"fmt"
)

type memoryService struct {
}

func newSMSMemory() Service {
	return &memoryService{}
}
func (M *memoryService) Send(ctx context.Context, message Message) error {
	fmt.Println(message)
	return nil
}
