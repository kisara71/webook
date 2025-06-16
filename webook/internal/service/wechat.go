package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"net/http"
)

const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
const authURLRedirect = "localhost:8080/oauth2/wechat/callback"
const verifyPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"

type WechatService interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error)
}

func NewWechatService(appId string) WechatService {
	return &wechatService{
		appId:  appId,
		secret: "124",
		client: http.DefaultClient,
	}
}

type wechatService struct {
	appId  string
	secret string // change to env
	client *http.Client
}

func (w *wechatService) VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error) {

	target := fmt.Sprintf(verifyPattern, w.appId, w.secret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)

	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	var res result
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, nil
	}

	return domain.WechatInfo{
		OpenID:  res.OpenID,
		UnionID: res.UnionID,
	}, nil
}

func (w *wechatService) AuthURL(ctx context.Context) (string, error) {
	state := uuid.New()
	return fmt.Sprintf(authURLPattern, w.appId, authURLRedirect, state), nil
}

type result struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int64  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}
