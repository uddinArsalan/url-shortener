package worker

import (
	"context"
	"log"
	"sync"
	"time"
	"url_shortener/internals/db"
	"url_shortener/models"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

func processClickQueue() {
	var rdb = db.GetRedisClient()
	ctx := context.Background()
	var wg sync.WaitGroup
	workerIDs := []string{"click_worker_1", "click_worker_2"}
	for _, workerID := range workerIDs {
		wg.Add(1)
		go func(workerID string) {
			defer wg.Done()
			for {
				streams := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    "click_worker",
					Consumer: workerID,
					Streams:  []string{"clicks:queue", ">"},
					Count:    10,
				})

				for _, stream := range streams.Val() {
					for _, message := range stream.Messages {
						values := message.Values
						analytics := models.ClickAnalytics{
							ID:        values["id"].(snowflake.ID),
							Timestamp: values["timestamp"].(time.Time),
							Ip:        values["ip"].(string),
							ShortCode: values["shortCode"].(string),
							Referer:   values["referer"].(string),
							Country:   values["country"].(string),
							City:      values["city"].(string),
							Os:        values["os"].(string),
							Browser:   values["browser"].(string),
							Device:    values["device"].(string),
						}
						err := db.InsertAnalyticsData(analytics)
						if err != nil {
							log.Printf("Error inserting analytics data: %v, Message ID: %v", err, message)
							continue
						}
						rdb.XAck(ctx, "clicks:queue", "click_worker", message.ID)
					}
				}
			}

		}(workerID)
	}
	wg.Wait()
	log.Println("Click queue processing stopped")
}

func main() {
	if err := db.InitRedis(); err != nil {
        log.Fatalf("Failed to initialize Redis: %v", err)
    }
	processClickQueue()
}
