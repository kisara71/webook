package sms

import (
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"os"
)

type aliService struct {
	client       *dysmsapi20170525.Client
	signName     string
	templateCode string
}

func newSMSAliYun(endpoint, signName, templateCode string) (Service, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
		Endpoint:        tea.String(endpoint),
	}
	client, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &aliService{
		client:       client,
		signName:     signName,
		templateCode: templateCode,
	}, nil
}

func (S *aliService) Send(ctx context.Context, msg Message) error {
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(msg.PhoneNumbers),

		SignName:      tea.String(S.signName),
		TemplateCode:  tea.String(S.templateCode),
		TemplateParam: tea.String(msg.TemplateParm),
	}
	_, err := S.client.SendSms(request)
	if err != nil {
		//	log
		return err
	}
	// log
	return nil
}
