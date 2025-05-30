package config

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	// "fmt"
	"strings"
	"url_shortener/internals/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
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

func GetClicksByDimension(w http.ResponseWriter, r *http.Request, dimension string) {
	vars := mux.Vars(r)
	urlId := vars["urlId"]
	if urlId == "" {
		http.Error(w, "URL ID required", http.StatusBadRequest)
		return
	}
	rdb := db.GetRedisClient()
	key := "clicks:" + urlId + ":" + dimension
	data, _ := rdb.ZRangeWithScores(r.Context(), key, 0, -1).Result()
	result := []map[string]interface{}{}
	label := strings.TrimPrefix(dimension, "by_")
	for _, z := range data {
		result = append(result, map[string]interface{}{
			label: z.Member.(string),
			"count": int(z.Score),
		})
	}
	json.NewEncoder(w).Encode(result)
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
