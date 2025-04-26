package db

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"os"
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

	fmt.Println("Connected to Redis successfully")
	return nil

}

func GetRedisClient() *redis.Client {
	return rdb
}
