package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"url_shortener/internals/db"
	"url_shortener/models"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

func generateRandString(nBytes int) (string, error) {
	bytes := make([]byte, nBytes)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name string, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func (kc *KeycloakAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	nonce, err := generateRandString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	setCallbackCookie(w, r, "state", state)
	setCallbackCookie(w, r, "nonce", nonce)
	http.Redirect(w, r, kc.Oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookies := []string{"token", "state", "nonce"}
	for _, cookie := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     cookie,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
	}
	// http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func (kc *KeycloakAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	verifier := kc.Provider.Verifier(kc.OIDCConfig)
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}
	oauth2Token, err := kc.Oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		http.Error(w, "nonce not found", http.StatusBadRequest)
		return
	}
	if idToken.Nonce != nonce.Value {
		http.Error(w, "nonce did not match", http.StatusBadRequest)
		return
	}

	oauth2Token.AccessToken = "*REDACTED*"

	var claims struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := db.FindUserByEmail(claims.Email)
	fmt.Println("User found :", user, "Error", err)
	if err == sql.ErrNoRows {
		user = models.User{
			Username: claims.Name,
			Email:    claims.Email,
		}
		err = db.InsertUser(user)
		if err != nil {
			log.Println("Failed to insert user", err)
			return
		}
	} else if err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtSecretKey) == 0 {
		http.Error(w, "JWT secret key not set", http.StatusInternalServerError)
		return
	}
	jwtClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := jwtToken.SignedString(jwtSecretKey)
	if err != nil {
		http.Error(w, "Failed to sign JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(24 * time.Hour.Seconds()),
	})
	redirectUrl := os.Getenv("REDIRECT_URL")
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
