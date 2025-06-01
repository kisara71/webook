package sms

type Message struct {
	PhoneNumbers string
	SignName     string
	TemplateCode string
	TemplateParm string
}

type Service interface {
	Send(msg Message) error
}
