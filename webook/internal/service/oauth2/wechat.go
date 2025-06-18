package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	oauth2Config "github.com/kisara71/WeBook/webook/config/oauth2"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"net/http"
)

//const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
//const authURLRedirect = "localhost:8080/oauth2/wechat/callback"
//const verifyPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"

func NewWechatService(config oauth2Config.Config) Service {
	return &wechatService{
		appId:             config.ClientID,
		secret:            config.ClientSecret,
		authURLPattern:    config.AuthURLPattern,
		redirectURL:       "http://localhost:8080/oauth2/wechat/callback",
		getBindURLPattern: config.ExchangeCodeURLPattern,
		client:            http.DefaultClient,
	}
}

type wechatService struct {
	appId             string
	secret            string // change to env
	authURLPattern    string
	getBindURLPattern string
	redirectURL       string
	client            *http.Client
}

func (w *wechatService) ExchangeCode(ctx context.Context, code string, state string) (domain.Oauth2Binding, error) {

	target := fmt.Sprintf(w.getBindURLPattern, w.appId, w.secret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)

	if err != nil {
		return domain.Oauth2Binding{}, err
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return domain.Oauth2Binding{}, err
	}

	var res wechatResult
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&res)
	if err != nil || res.ErrCode != 0 {
		return domain.Oauth2Binding{}, nil
	}

	return domain.Oauth2Binding{

		ExternalID: res.OpenID,
		Provider:   domain.ProviderWechat,
	}, nil
}

func (w *wechatService) AuthURL(ctx context.Context) string {
	state := uuid.New()
	return fmt.Sprintf(w.authURLPattern, w.appId, w.redirectURL, state)
}

type wechatResult struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int64  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}
