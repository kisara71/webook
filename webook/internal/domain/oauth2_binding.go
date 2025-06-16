package domain

type ProviderType string

const (
	ProviderUnknown = "UNKNOWN"
	ProviderWechat  = "WECHAT"
	ProviderGithub  = "GITHUB"
)

type Oauth2Binding struct {
	UserID int64
	ID     int64

	Provider   ProviderType
	ExternalID string
}

func (p ProviderType) ToString() string {
	return string(p)
}
