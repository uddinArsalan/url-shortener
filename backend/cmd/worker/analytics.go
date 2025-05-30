package worker

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"
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
						idStr := values["ID"].(string)
						id, err := strconv.ParseInt(idStr, 10, 64)
						if err != nil {
							log.Printf("Error parsing ID: %v", err)
							continue
						}

						timestampStr := values["Timestamp"].(string)
						timestamp, err := time.Parse(time.RFC3339, timestampStr)
						if err != nil {
							log.Printf("Error parsing timestamp: %v", err)
							continue
						}
						analytics := models.ClickAnalytics{
							ID:        id,
							Timestamp: timestamp,
							Ip:        values["Ip"].(string),
							ShortCode: values["ShortCode"].(string),
							Referrer:  values["Referrer"].(string),
							Country:   values["Country"].(string),
							City:      values["City"].(string),
							Os:        values["Os"].(string),
							Browser:   values["Browser"].(string),
							Device:    values["Device"].(string),
						}
						err = db.InsertAnalyticsData(analytics)
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
