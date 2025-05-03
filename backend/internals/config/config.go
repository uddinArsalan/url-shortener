package config

import (
	"context"
	"net/http"
	"os"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type KeycloakConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Realm        string
	BaseURL      string
}

type contextKey string

const UserContextKey contextKey = "user"

type HTTPError struct  {
	Status int
	Message string
}

func GetUserIDFromContext(ctx context.Context) (int64, *HTTPError) {
	claims, ok := ctx.Value(UserContextKey).(jwt.MapClaims)
	if !ok {
		return 0, &HTTPError{Status: http.StatusUnauthorized, Message: "Unauthorized"}
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, &HTTPError{Status: http.StatusBadRequest, Message: "Invalid user ID"}
	}

	return int64(sub), nil
}


func LoadKeycloakConfig() KeycloakConfig {
	return KeycloakConfig{
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("KEYCLOAK_REDIRECT_URL"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
		BaseURL:      os.Getenv("KEYCLOAK_BASE_URL"),
	}
}
