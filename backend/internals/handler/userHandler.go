package handler

import (
	"encoding/json"
	"net/http"
	"url_shortener/internals/config"
	"url_shortener/internals/db"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID, er := config.GetUserIDFromContext(r.Context())
	if er != nil {
		http.Error(w, er.Message, er.Status)
		return
	}
	user, err := db.FindUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
