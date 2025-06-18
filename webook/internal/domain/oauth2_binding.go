package domain

const (
	ProviderUnknown = "UNKNOWN"
	ProviderWechat  = "WECHAT"
	ProviderGithub  = "GITHUB"
)

type Oauth2Binding struct {
	UserID int64
	ID     int64

	Provider   string
	ExternalID string

	AccessToken string
}
