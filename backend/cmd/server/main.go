package server

import (
	"context"
	"log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"net/http"

	// "github.com/syumai/workers"
	"url_shortener/cmd/worker"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kcAuth, err := auth.InitKeycloak(ctx, cfg)
	if err != nil {
		panic(err)
	}
	go worker.ProcessClickQueue(ctx)

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	public := apiRouter.NewRoute().Subrouter()
	public.HandleFunc("/auth/login", kcAuth.HandleLogin).Methods("GET")
	public.HandleFunc("/auth/callback", kcAuth.HandleCallback).Methods("GET")
	protected := apiRouter.NewRoute().Subrouter()
	public.Use(middleware.PerClientRateLimiter)
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/shorten", handler.ShortenURL).Methods("POST")
	protected.HandleFunc("/auth/logout", auth.HandleLogout).Methods("GET")
	public.HandleFunc("/url/{shortCode}", func(w http.ResponseWriter, r *http.Request) {
		middleware.TrackClickMiddleware(http.HandlerFunc(handler.RedirectURL)).ServeHTTP(w, r)
	}).Methods("GET")
	protected.HandleFunc("/me", handler.MeHandler).Methods("GET")
	protected.HandleFunc("/url", handler.GetUserUrls).Methods("GET")
	protected.HandleFunc("/analytics", handler.AnalyticsOfURL).Methods("GET")
	protected.HandleFunc("/analytics/{urlId}/hourly", handler.GetHourlyClicks).Methods("GET")
	protected.HandleFunc("/analytics/{urlId}/country", handler.GetCountryWiseClicks).Methods("GET")
	protected.HandleFunc("/analytics/{urlId}/city", handler.GetCityWiseClicks).Methods("GET")
	protected.HandleFunc("/analytics/{urlId}/device", handler.GetDeviceWiseClicks).Methods("GET")
	protected.HandleFunc("/analytics/{urlId}/browser", handler.GetBrowserWiseClicks).Methods("GET")
	protected.HandleFunc("/analytics/{urlId}/referrer", handler.GetReferrerWiseClicks).Methods("GET")
	handlerWithCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", ""},
		AllowCredentials: true,
	}).Handler(r)

	// workers.Serve(handlerWithCors)

	if err := http.ListenAndServe(":4000", handlerWithCors); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
