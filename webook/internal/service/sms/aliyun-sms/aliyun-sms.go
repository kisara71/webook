package aliyun_sms

import (
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/kisara71/WeBook/webook/internal/service/sms"
	"os"
)

type Service struct {
	client       *dysmsapi20170525.Client
	signName     string
	templateCode string
}

func NewService(endpoint, signame, templateCode string) (*Service, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
		Endpoint:        tea.String(endpoint),
	}
	client, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Service{
		client:       client,
		signName:     signame,
		templateCode: templateCode,
	}, nil
}

func (s *Service) Send(ctx context.Context, msg sms.Message) error {
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(msg.PhoneNumbers),
		SignName:      tea.String(s.signName),
		TemplateCode:  tea.String(s.templateCode),
		TemplateParam: tea.String(msg.TemplateParm),
	}
	_, err := s.client.SendSms(request)
	if err != nil {
		//	log
		return err
	}
	// log
	return nil
}
