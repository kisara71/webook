package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	oauth2Config "github.com/kisara71/WeBook/webook/config/oauth2"
	"github.com/kisara71/WeBook/webook/internal/domain"
	"net/http"
	"net/url"
	"strconv"
)

type githubService struct {
	clientID        string
	clientSecret    string
	authURLPattern  string
	redirectURL     string
	exchangeCodeURL string
	userInfoURL     string
	scope           string
	client          *http.Client
}

func NewGithubService(config oauth2Config.Config) Service {
	proxy, _ := url.Parse("http://127.0.0.1:7897")
	return &githubService{
		clientID:        config.ClientID,
		clientSecret:    config.ClientSecret,
		redirectURL:     "http://127.0.0.1:8080/oauth2/github/callback",
		authURLPattern:  config.AuthURLPattern,
		exchangeCodeURL: config.ExchangeCodeURLPattern,
		userInfoURL:     "https://api.github.com/user",
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		},
		//scope: config.Scope,
	}
}

func (g *githubService) AuthURL(ctx context.Context) string {
	if g.clientID == "" || g.authURLPattern == "" {
		return ""
	}
	state := uuid.New().String()
	return fmt.Sprintf(g.authURLPattern, g.clientID, g.redirectURL, state)
}

func (g *githubService) ExchangeCode(ctx context.Context, code string, state string) (domain.Oauth2Binding, error) {
	type reqBody struct {
		Code         string `json:"code"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURL  string `json:"redirect_url"`
	}
	data := reqBody{
		Code:         code,
		ClientSecret: g.clientSecret,
		ClientID:     g.clientID,
		RedirectURL:  g.redirectURL,
	}
	body, err := json.Marshal(&data)
	if err != nil {
		return domain.Oauth2Binding{}, fmt.Errorf("marshal request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.exchangeCodeURL, bytes.NewBuffer(body))
	if err != nil {
		return domain.Oauth2Binding{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return domain.Oauth2Binding{}, fmt.Errorf("do request: %w", err)
	}

	var result githubExchangeCodeResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return domain.Oauth2Binding{}, fmt.Errorf("decode response: %w", err)
	}

	if result.AccessToken == "" {
		return domain.Oauth2Binding{}, fmt.Errorf("empty access token")
	}
	userinfo, err := g.getUserInfo(ctx, result.AccessToken)
	if err != nil {
		return domain.Oauth2Binding{}, err
	}
	return domain.Oauth2Binding{
		ExternalID:  strconv.FormatInt(userinfo.ID, 10),
		Provider:    domain.ProviderGithub,
		AccessToken: result.AccessToken,
	}, nil
}

func (g *githubService) getUserInfo(ctx context.Context, acToken string) (githubGetUserResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.userInfoURL, nil)
	if err != nil {
		return githubGetUserResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+acToken)
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return githubGetUserResult{}, err
	}
	var res githubGetUserResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		return githubGetUserResult{}, err
	}
	
	return res, nil
}

type githubExchangeCodeResult struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

type githubGetUserResult struct {
	ID int64 `json:"id"`
}
