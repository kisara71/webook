package sms

import "context"

type Message struct {
	PhoneNumbers string
	TemplateParm string
}

type Service interface {
	Send(context.Context, Message) error
}

func NewSMSService() Service {
	return newSMSMemory()
}
