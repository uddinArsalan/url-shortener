package auth

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
	"url_shortener/internals/config"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type KeycloakAuth struct {
	Provider     *oidc.Provider
	Oauth2Config oauth2.Config
	OIDCConfig   *oidc.Config
	Config       config.KeycloakConfig
}

// func InitKeycloak(ctx context.Context, cfg config.KeycloakConfig) (*KeycloakAuth, error) {
// 	encodedRealm := url.PathEscape(cfg.Realm)
// 	issuer := fmt.Sprintf("%s/realms/%s", cfg.BaseURL, encodedRealm)

// 	provider, err := oidc.NewProvider(ctx, issuer)
// 	if err != nil {
// 		fmt.Printf("Error %v\n", err)
// 		return nil, err
// 	}
// 	oidcConfig := &oidc.Config{
// 		ClientID: cfg.ClientID,
// 	}
// 	oauth2Config := oauth2.Config{
// 		ClientID:     cfg.ClientID,
// 		ClientSecret: cfg.ClientSecret,
// 		RedirectURL:  cfg.RedirectURL,

// 		Endpoint: provider.Endpoint(),
// 		Scopes:   []string{oidc.ScopeOpenID, "profile", "email"},
// 	}
// 	return &KeycloakAuth{
// 		Provider:     provider,
// 		Oauth2Config: oauth2Config,
// 		OIDCConfig:   oidcConfig,
// 	}, nil
// }

// quick fix will checkout again cause or issue
// func InitKeycloak(ctx context.Context, cfg config.KeycloakConfig) (*KeycloakAuth, error) {
// 	encodedRealm := url.PathEscape(cfg.Realm)
// 	issuer := fmt.Sprintf("%s/realms/%s", cfg.BaseURL, encodedRealm)

// 	var provider *oidc.Provider
// 	var err error

// 	for i := 1; i <= 5; i++ {
// 		fmt.Printf("Keycloak discovery attempt %d: %s\n", i, issuer)

// 		provider, err = oidc.NewProvider(ctx, issuer)

// 		if err == nil {
// 			break
// 		}

// 		if strings.Contains(err.Error(), "429") {
// 			wait := time.Duration(i*2) * time.Second
// 			fmt.Printf("Got 429 from Keycloak. Retrying in %v...\n", wait)
// 			time.Sleep(wait)
// 			continue
// 		}
// 		break
// 	}

// 	if err != nil {
// 		fmt.Printf("WARNING: Keycloak initialization failed: %v\n", err)
// 		return nil, nil
// 	}

// 	oidcConfig := &oidc.Config{
// 		ClientID: cfg.ClientID,
// 	}

// 	oauth2Config := oauth2.Config{
// 		ClientID:     cfg.ClientID,
// 		ClientSecret: cfg.ClientSecret,
// 		RedirectURL:  cfg.RedirectURL,
// 		Endpoint:     provider.Endpoint(),
// 		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
// 	}

// 	return &KeycloakAuth{
// 		Provider:     provider,
// 		Oauth2Config: oauth2Config,
// 		OIDCConfig:   oidcConfig,
// 	}, nil
// }

func InitKeycloak(ctx context.Context, cfg config.KeycloakConfig) (*KeycloakAuth, error) {
	encodedRealm := url.PathEscape(cfg.Realm)
	issuer := fmt.Sprintf("%s/realms/%s", cfg.BaseURL, encodedRealm)

	provider, err := tryDiscovery(ctx, issuer)
	if err != nil {
		// Log a warning but return a partially initialized KeycloakAuth
		// Provider is nil, discovery will be retried lazily
		fmt.Printf("Keycloak initialization failed at startup: %v\n", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		// Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	
	if provider != nil {
		oauth2Config.Endpoint = provider.Endpoint()
	}

	ka := &KeycloakAuth{
		Provider:     provider,
		Oauth2Config: oauth2Config,
		OIDCConfig:   oidcConfig,
		Config:       cfg,
	}

	return ka, nil
}

func tryDiscovery(ctx context.Context, issuer string) (*oidc.Provider, error) {
	var provider *oidc.Provider
	var err error

	for i := 1; i <= 5; i++ {
		fmt.Printf("Keycloak discovery attempt %d: %s\n", i, issuer)
		provider, err = oidc.NewProvider(ctx, issuer)
		if err == nil {
			return provider, nil
		}

		if strings.Contains(err.Error(), "429") {
			wait := time.Duration(i*2) * time.Second
			fmt.Printf("Got 429 from Keycloak. Retrying in %v...\n", wait)
			time.Sleep(wait)
			continue
		}

		break
	}

	return nil, err
}
