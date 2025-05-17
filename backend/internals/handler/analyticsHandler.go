package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"url_shortener/internals/db"
)

func AnalyticsOfURL(w http.ResponseWriter, r *http.Request) {
	urlId := r.URL.Query().Get("urlId")
	if urlId == "" {
		http.Error(w, "URL ID required", http.StatusBadRequest)
		return
	}
	urlID, convErr := strconv.ParseInt(urlId, 10, 64)
	if convErr != nil {
		http.Error(w, "Invalid urlId", http.StatusBadRequest)
		return
	}
	userAnalytics, err := db.FindUserAnaltics(urlID)
	if err != nil {
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userAnalytics)
}
