package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"url_shortener/internals/config"
	"url_shortener/internals/db"
	"url_shortener/models"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func generateShortCode() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[r.Intn(len(chars))]
	}
	return string(code)
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	userUrl := r.URL.Query().Get("url")
	if userUrl == "" {
		http.Error(w, "URL required", http.StatusBadRequest)
		return
	}
	userID, er := config.GetUserIDFromContext(r.Context())
	if er != nil {
		http.Error(w, er.Message, er.Status)
		return
	}
	shortCode := generateShortCode()
	rdb := db.GetRedisClient()
	ctx := context.Background()
	if err := rdb.Set(ctx, shortCode, userUrl, 0).Err(); err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}
	url := models.URL{
		OriginalURL: userUrl,
		ShortCode:   shortCode,
		UserID:      userID,
	}
	if err := db.InsertUrl(url); err != nil {
		http.Error(w, "Failed to save URL to database", http.StatusInternalServerError)
		return
	}
	w.Write(fmt.Appendf(nil, "http://localhost:4000/api/v1/url/%s", shortCode))
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shortCode := vars["shortCode"]
	rdb := db.GetRedisClient()
	ctx := context.Background()
	url, err := rdb.Get(ctx, shortCode).Result()
	if err != nil {
		if err == redis.Nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Failed to retrieve URL", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}


func GetUserUrls(w http.ResponseWriter, r *http.Request) {
	userId , err := config.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Message, err.Status)
		return
	}
	queryParams := r.URL.Query()
    cursor := queryParams.Get("cursor")
    limit := queryParams.Get("limit")
	
	if limit == "" {
		limit = "3"
	}
	limitInt, convErr := strconv.Atoi(limit)
	if convErr != nil {	
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	urls, dbErr := db.FindUrlsFromUserId(strconv.FormatInt(userId, 10),limitInt,cursor)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)

}