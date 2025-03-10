package db

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return UserService{db: db}
}

func NewClient() {
	DB_NAME := os.Getenv("DB_NAME")
	DB_HOST := os.Getenv("DB_HOST")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	dbUri := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", DB_USERNAME, DB_PASSWORD, DB_HOST, DB_NAME)

	db, err := sql.Open("mysql", dbUri)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("Database connected successfully!")

	service := NewUserService(db)
	fmt.Printf("Service is ready! %v\n", service)

	createUserTable(service.db)
	InsertUser(User{"Arsu", "arsalan91@gmail.com", "random_password"}, service.db)
}

func createUserTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users(
	id BIGINT PRIMARY KEY ,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	data, err := db.Exec(query)

	if err == nil {
		fmt.Println("Users Table Created ", data)
	} else {
		log.Fatalf("Error creating users table %v", err)
	}

}

type User struct {
	username      string
	email         string
	password_hash string
}

func InsertUser(user User, db *sql.DB) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := node.Generate()
	query := `INSERT INTO users (id,username, email,password_hash, created_at) VALUES (?, ?, ?,?,?)`
	result, err := db.Exec(query, id, user.username, user.email, user.password_hash, time.Now())
	if err != nil {
		log.Fatalf("Error inserting data in users table %v", err)
	} else {
		fmt.Println("User inserted successfully!")
	}
	fmt.Println(result.LastInsertId())
}
