package worker

import (
	"context"
	"log"
	"sync"
	"url_shortener/internals/db"
	"url_shortener/models"

	"github.com/redis/go-redis/v9"
)

func ProcessClickQueue(ctx context.Context) {
	var rdb = db.GetRedisClient()
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
							ID:        values["ID"].(string),
							Timestamp: values["Timestamp"].(string),
							Ip:        values["Ip"].(string),
							ShortCode: values["ShortCode"].(string),
							Referer:   values["Referer"].(string),
							Country:   values["Country"].(string),
							City:      values["City"].(string),
							Os:        values["Os"].(string),
							Browser:   values["Browser"].(string),
							Device:    values["Device"].(string),
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
