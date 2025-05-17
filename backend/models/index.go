package models

import (
	"time"
)

type User struct {
	ID        int64  `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type URL struct {
	ID          string     `json:"id" db:"id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	ShortCode   string    `json:"shortcode" db:"shortcode"`
	UserID      int64     `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Pagination struct {
	NextCursor string `json:"next_cursor"`
	HasMore bool `json:"has_more"`
}

type URLResponse struct {
	Urls []URL `json:"urls"`
	Pagination Pagination `json:"pagination"`
}

type ClickAnalytics struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Ip        string    `json:"ip_hash"`
	ShortCode string    `json:"shortCode"`
	Referrer   string    `json:"referrer"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Os        string    `json:"os"`
	Browser   string    `json:"browser"`
	Device    string    `json:"device"`
}

type UserAnalytics struct {
	ClickAnalytics []ClickAnalytics `json:"click_analytics"`
	TotalClicks int64 `json:"total_clicks"`
	UniqueClicks int64 `json:"unique_clicks"`
}
