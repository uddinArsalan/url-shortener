package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"math/rand"
	"time"
	// "rate_limiter/db"
)

var urlMap = make(map[string]string)

func generateShortCode() string{
	rand.New(rand.NewSource(time.Now().Unix()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte,6)
	for i := range(code){
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func shorten(w http.ResponseWriter, r *http.Request){
	userUrl := r.URL.Query().Get("url")
	if userUrl == ""{
		http.Error(w, "URL required", http.StatusBadRequest)
		fmt.Println("No url provided to shorten")
		return
	}
	shortCode := generateShortCode()
	urlMap[shortCode] = userUrl
	fmt.Printf("\nUrl shortened from %s to %s",userUrl,shortCode)
}

// func redirect(w http.ResponseWriter, r *http.Request){
// 	urlShortCode := r.URL.Path
// 	fmt.Println(urlShortCode)
// }

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/shorten", shorten).Methods("POST")

	r.HandleFunc("/{shortCode}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortCode := vars["shortCode"]
		url,exists := urlMap[shortCode]
		if !exists {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}).Methods("GET")
	
	// db.NewClient()
	
	http.ListenAndServe(":80", r)
}
