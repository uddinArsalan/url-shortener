package main

import (
	"context"
	"log"
	"url_shortener/internals/db"
)

func main() {
	ctx := context.Background()

	db.InitDBClient()
	if err := db.InitRedis(); err != nil {
		log.Fatalf("Redis init failed: %v", err)
	}

	ProcessClickQueue(ctx)
}
