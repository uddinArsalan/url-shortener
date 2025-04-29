package handler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"url_shortener/internals/db"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func generateShortCode() string {
	rand.New(rand.NewSource(time.Now().Unix()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	userUrl := r.URL.Query().Get("url")
	if userUrl == "" {
		http.Error(w, "URL required", http.StatusBadRequest)
		fmt.Println("No url provided to shorten")
		return
	}
	shortCode := generateShortCode()
	rdb := db.GetRedisClient()
	ctx := context.Background()
	if err := rdb.Set(ctx, shortCode, userUrl, 0).Err(); err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}
	url := db.URL{
		OriginalUrl: userUrl,
		ShortCode:   shortCode,
		UserId:      1902446527184375808,
	}
	db.InsertUrl(url)
	fmt.Printf("\nUrl shortened from %s to %s", userUrl, shortCode)
	w.Write(fmt.Appendf(nil, "http://localhost/url/%s", shortCode))
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