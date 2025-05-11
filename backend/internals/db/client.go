package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"url_shortener/models"

	"github.com/bwmarrin/snowflake"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDBClient() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("sql.Open failed: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatalf("Database connection failed: %v", err)
	}
	fmt.Println("Database connected successfully!")

	CreateUserTable()
	CreateUrlTable()
	CreateAnalyticsTable()
}

func CreateUserTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users(
	id BIGINT PRIMARY KEY ,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)
	`
	data, err := db.Exec(query)

	if err == nil {
		rowsAffected, err := data.RowsAffected()
		if err != nil {
			log.Fatalf("Error fetching rows affected: %v", err)
		}
		fmt.Printf("Users Table Created, Rows Affected: %d\n", rowsAffected)
	} else {
		log.Fatalf("Error creating users table %v", err)
	}

}

func CreateUrlTable() {
	query := `
	CREATE TABLE IF NOT EXISTS urls(
	id BIGINT PRIMARY KEY,
	original_url TEXT NOT NULL,
	shortcode VARCHAR(7) UNIQUE NOT NULL,
	user_id BIGINT references users(id),
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)
	`
	data, err := db.Exec(query)
	if err == nil {
		rowsAffected, err := data.RowsAffected()
		if err != nil {
			log.Fatalf("Error fetching rows affected: %v", err)
		}
		fmt.Printf("Url Table Created, Rows Affected: %d\n", rowsAffected)
	} else {
		log.Fatalf("Error creating url table %v", err)
	}
}

func CreateAnalyticsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS analytics(
	id BIGINT PRIMARY KEY,
	url_id BIGINT REFERENCES urls(id), 
	ip_hash VARCHAR(64) NOT NULL,
    country VARCHAR(50),
    city VARCHAR(100),
    os VARCHAR(50),
    browser VARCHAR(50),
    device VARCHAR(20),
    referrer TEXT,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)
	`
	data, err := db.Exec(query)
	if err == nil {
		rowsAffected, err := data.RowsAffected()
		if err != nil {
			log.Fatalf("Error fetching rows affected: %v", err)
		}
		fmt.Printf("Analytics Table Created, Rows Affected: %d\n", rowsAffected)
	} else {
		log.Fatalf("Error creating Analytics table %v", err)
	}
}

func FindUrlIdFromShortCode(shortCode string) (int64, error) {
	query := `SELECT id FROM urls WHERE shortcode = $1`
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var urlId int64
	er := stmt.QueryRow(shortCode).Scan(&urlId)
	if er != nil {
		return 0, err
	}
	return urlId, nil
}

func InsertUser(user models.User) error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return fmt.Errorf("failed to create snowflake node: %w", err)
	}
	id := node.Generate()
	query := `INSERT INTO users (id,username, email, created_at) VALUES ($1, $2, $3, $4)`
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id.Int64(), user.Username, user.Email, time.Now())
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}
	fmt.Println("User inserted successfully!")
	return nil
}

func FindUserByEmail(email string) (models.User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE email = $1`
	stmt, err := db.Prepare(query)
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()
	var user models.User

	err = stmt.QueryRow(email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func FindUserByID(id int64) (models.User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	stmt, err := db.Prepare(query)
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()
	var user models.User

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("no user found with id: %d", id)
		}
		return models.User{}, fmt.Errorf("error querying user: %w", err)
	}
	return user, nil
}

func InsertUrl(url models.URL) error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return fmt.Errorf("failed to create snowflake node: %w", err)
	}
	id := node.Generate()
	query := `INSERT INTO urls (id,original_url,shortcode,user_id, created_at) VALUES ($1, $2, $3, $4, $5)`
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id.Int64(), url.OriginalUrl, url.ShortCode, url.UserId, time.Now())
	if err != nil {
		log.Fatalf("Error inserting data in urls table %v", err)
		return err
	}
	fmt.Println("Urls inserted successfully!")
	return nil
}

func InsertAnalyticsData(clicksData models.ClickAnalytics) error {
	query := `INSERT INTO analytics (id,url_id,ip_hash,country,city,os,browser,device,referrer,timestamp) VALUES ($1, $2, $3, $4, $5,$6,$7,$8,$9,$10)`
	urlId, err := FindUrlIdFromShortCode(clicksData.ShortCode)
	if urlId == 0 || err != nil {
		log.Printf("Error finding URL ID for shortcode %s: %v", clicksData.ShortCode, err)
		return fmt.Errorf("error finding URL ID for shortcode %s: %w", clicksData.ShortCode, err)
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()
	id, err := strconv.ParseInt(clicksData.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse clicksData.ID to int64: %w", err)
	}
	_, err = stmt.Exec(id, urlId, clicksData.Ip, clicksData.Country, clicksData.City, clicksData.Os, clicksData.Browser, clicksData.Device, clicksData.Referer, clicksData.Timestamp)
	if err != nil {
		log.Fatalf("Error inserting data in analytics table %v", err)
		return err
	}
	return nil
}
