package handler

import (
	"encoding/json"
	"net/http"
	// "strconv"
	"time"
	"url_shortener/internals/config"
	"url_shortener/internals/db"

	"github.com/gorilla/mux"
)

func AnalyticsOfURL(w http.ResponseWriter, r *http.Request) {
	urlId := r.URL.Query().Get("urlId")
	if urlId == "" {
		http.Error(w, "URL ID required", http.StatusBadRequest)
		return
	}
	// urlID, convErr := strconv.ParseInt(urlId, 10, 64)
	// if convErr != nil {
	// 	http.Error(w, "Invalid urlId", http.StatusBadRequest)
	// 	return
	// }
	userAnalytics, err := db.FindUserAnaltics(urlId)
	if err != nil {
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userAnalytics); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetHourlyClicks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlId := vars["urlId"]
	if urlId == "" {
		http.Error(w, "URL ID required", http.StatusBadRequest)
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	rdb := db.GetRedisClient()

	layout := time.RFC3339
	fromTime, _ := time.Parse(layout, fromStr)
	toTime, _ := time.Parse(layout, toStr)

	data, _ := rdb.ZRangeWithScores(r.Context(), "clicks:"+urlId+":by_hour", 0, -1).Result()
	result := []map[string]interface{}{}

	for _, z := range data {
		t, err := time.Parse(layout, z.Member.(string))
		// fmt.Printf("Time %v, FromTime%v, ToTime %v", t, fromTime, toTime)
		// fmt.Println(err)
		if err != nil {
			continue
		}
		if t.After(fromTime) && t.Before(toTime) {
			// fmt.Println("Inside Time ")
			result = append(result, map[string]interface{}{
				"hour":  z.Member.(string),
				"count": int(z.Score),
			})
		}
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetCountryWiseClicks(w http.ResponseWriter, r *http.Request) {
	config.GetClicksByDimension(w, r, "by_country")
}

func GetCityWiseClicks(w http.ResponseWriter, r *http.Request) {
	config.GetClicksByDimension(w, r, "by_city")
}

func GetDeviceWiseClicks(w http.ResponseWriter, r *http.Request) {
	config.GetClicksByDimension(w, r, "by_device")
}

func GetBrowserWiseClicks(w http.ResponseWriter, r *http.Request) {
	config.GetClicksByDimension(w, r, "by_browser")
}

func GetReferrerWiseClicks(w http.ResponseWriter, r *http.Request) {
	config.GetClicksByDimension(w, r, "by_referrer")
}
