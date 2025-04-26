package config

import (
	"os"
	_ "github.com/joho/godotenv/autoload"
)

type KeycloakConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Realm        string
	BaseURL      string
}

func LoadKeycloakConfig() KeycloakConfig {
	return KeycloakConfig{
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("KEYCLOAK_REDIRECT_URL"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
		BaseURL:      "http://localhost:8080",
	}
}
