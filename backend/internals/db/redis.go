package db

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() error {
	ctx := context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	err = rdb.XGroupCreateMkStream(ctx,"clicks:queue","click_worker","0").Err()
	if err != nil {
		log.Printf("Error creating Redis stream group: %v", err)
	}

	fmt.Println("Connected to Redis successfully")
	return nil

}

func GetRedisClient() *redis.Client {
	return rdb
}
