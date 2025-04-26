package auth

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"net/url"
	"url_shortener/internals/config"
)

type KeycloakAuth struct {
	Provider     *oidc.Provider
	Oauth2Config oauth2.Config
	OIDCConfig   *oidc.Config
}

func InitKeycloak(ctx context.Context, cfg config.KeycloakConfig) (*KeycloakAuth, error) {
	encodedRealm := url.PathEscape(cfg.Realm)
	issuer := fmt.Sprintf("%s/realms/%s", cfg.BaseURL, encodedRealm)

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}
	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,

		Endpoint: provider.Endpoint(),
		Scopes:   []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return &KeycloakAuth{
		Provider:     provider,
		Oauth2Config: oauth2Config,
		OIDCConfig:   oidcConfig,
	}, nil
}
