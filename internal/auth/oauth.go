package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Client struct {
	Provider *oidc.Provider
	OIDC     *oidc.IDTokenVerifier
	Oauth    oauth2.Config
}
type Config struct {
	BaseURL      string // Authorization base url
	ClientID     string // client id oauth
	RedirectURI  string // valid redirect uri
	ClientSecret string // optional
}

func New(ctx context.Context, config *Config) (*Client, error) {
	// Construct the provider URI oauth
	providerURL := config.BaseURL
	// Google's OAuth 2.0 endpoint
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}
	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})
	// Configure an OpenID Connect aware OAuth2 client
	oauth2 := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID,
			"profile",
			"email",
		},
	}
	return &Client{
		Oauth:    oauth2,
		OIDC:     verifier,
		Provider: provider,
	}, nil
}
