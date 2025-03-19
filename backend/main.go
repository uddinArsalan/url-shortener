package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"log"
	"math/rand"
	"net/http"
	"time"
	"url_shortener/db"
)

// var urlMap = make(map[string]string)

func generateShortCode() string {
	rand.New(rand.NewSource(time.Now().Unix()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func shorten(w http.ResponseWriter, r *http.Request) {
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
func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shortCode := vars["shortCode"]
	// url,exists := urlMap[shortCode]
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

func getAllUrls(w http.ResponseWriter, r *http.Request) {
	rdb := db.GetRedisClient()
	ctx := context.Background()
	urls := rdb.Scan(ctx, 0, "*", 3)
	fmt.Println(urls.Args()...)
	w.Write(fmt.Appendf(nil, "All Urls are http://localhost/url/%s", urls))
}

func main() {
	db.InitDBClient()
	if err := db.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/shorten", shorten).Methods("POST")

	r.HandleFunc("/url/{shortCode}", redirect).Methods("GET")

	r.HandleFunc("/get", getAllUrls).Methods("GET")
	r.HandleFunc("/createUser", func(w http.ResponseWriter, r *http.Request) {
		user := db.User{
			Username:     "Arsu",
			Email:        "arsu91@gmail.com",
			PasswordHash: "12364643",
		}
		db.InsertUser(user)
	}).Methods("POST")
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		user, err := db.FindUserByEmail(email)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
	})

	handler := cors.Default().Handler(r)

	http.ListenAndServe(":80", handler)
}
