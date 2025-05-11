package models

type User struct {
	ID        int64
	Username  string
	Email     string
	CreatedAt string
}

type URL struct {
	ID          int64
	OriginalUrl string
	ShortCode   string
	UserId      int64
	CreatedAt   string
}

type ClickAnalytics struct {
	ID        string `json:"id"`
	Timestamp string    `json:"timestamp"`
	Ip        string       `json:"ip"`
	ShortCode string       `json:"shortCode"`
	Referer   string       `json:"referer"`
	Country   string       `json:"country"`
	City      string       `json:"city"`
	Os        string       `json:"os"`
	Browser   string       `json:"browser"`
	Device    string       `json:"device"`
}
