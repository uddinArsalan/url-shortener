package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"net/http"

	"url_shortener/internals/auth"
	"url_shortener/internals/config"
	"url_shortener/internals/db"
	"url_shortener/internals/handler"
	"url_shortener/internals/middleware"
)

func Start() {
	db.InitDBClient()
	if err := db.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	cfg := config.LoadKeycloakConfig()
	ctx := context.Background()

	kcAuth, err := auth.InitKeycloak(ctx, cfg)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware.PerClientRateLimiter(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				next.ServeHTTP(writer, request)
			})).ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/shorten", handler.ShortenURL).Methods("POST")
	r.HandleFunc("/url/{shortCode}", handler.RedirectURL).Methods("GET")
	r.HandleFunc("/auth/login", kcAuth.HandleLogin).Methods("GET")
	r.HandleFunc("/auth/callback", kcAuth.HandleCallback).Methods("GET")
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		user, err := db.FindUserByEmail(email)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
	})

	handler := cors.Default().Handler(r)

	http.ListenAndServe(":4000", handler)
}
