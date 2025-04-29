package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var db *sql.DB

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
	Clicks      int32
}

func InitDBClient() {
	// DATABASE_NAME := os.Getenv("DATABASE_NAME")
	// DATABASE_HOST := os.Getenv("DATABASE_HOST")
	// DATABASE_USER := os.Getenv("DATABASE_USER")
	// DATABASE_PASSWORD := os.Getenv("DATABASE_PASSWORD")

	// connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname=%s", DATABASE_USER, DATABASE_PASSWORD, DATABASE_HOST, DATABASE_NAME)
	// dbUri := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", DB_USERNAME, DB_PASSWORD, DB_HOST, DB_NAME)
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    clicks INT DEFAULT 0
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
	ip_address VARCHAR(45),
    user_agent TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

func FindUrlFromShortCode(shortCode string) (string, error) {
	query := `SELECT original_url FROM urls WHERE shortcode = $1`
	stmt, err := db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var originalURL string
	er := stmt.QueryRow(shortCode).Scan(&originalURL)
	if er != nil {
		fmt.Printf("Error querying URL with given short code %v\n", er)
		if er == sql.ErrNoRows {
			return "", fmt.Errorf("no URL found for the given short code: %v", shortCode)
		}
		return "", err
	}
	return originalURL, nil
}

func InsertUser(user User) error {
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

func FindUserByEmail(email string) (User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE email = $1`
	stmt, err := db.Prepare(query)
	if err != nil {
		return User{}, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()
	var user User

	err = stmt.QueryRow(email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("no user found with email: %s", email)
		}
		return User{}, fmt.Errorf("error querying user: %w", err)
	}
	return user, nil
}

func FindUserByID(id int64) (User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	stmt, err := db.Prepare(query)
	if err != nil {
		return User{}, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()
	var user User

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("no user found with id: %d", id)
		}
		return User{}, fmt.Errorf("error querying user: %w", err)
	}
	return user, nil
}

func InsertUrl(url URL) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := node.Generate()
	// 1902446527184375808
	query := `INSERT INTO urls (id,original_url,shortcode,user_id, created_at,clicks) VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id.Int64(), url.OriginalUrl, url.ShortCode, url.UserId, time.Now(), url.Clicks)
	if err != nil {
		log.Fatalf("Error inserting data in urls table %v", err)
	} else {
		fmt.Println("Urls inserted successfully!")
	}
	fmt.Println(result)
}
