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
	ID          int64     `json:"id" db:"id"`
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
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Ip        string `json:"ip"`
	ShortCode string `json:"shortCode"`
	Referer   string `json:"referer"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Os        string `json:"os"`
	Browser   string `json:"browser"`
	Device    string `json:"device"`
}
