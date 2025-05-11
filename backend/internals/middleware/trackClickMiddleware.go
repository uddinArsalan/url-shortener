package middleware

import (
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"time"
	"url_shortener/internals/db"
	"url_shortener/models"

	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/mux"
	"github.com/mssola/useragent"
	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
)

func hashIp(ip string) string {
	h := sha256.New()
	h.Write([]byte(ip))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func TrackClickMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var rdb = db.GetRedisClient()
		db, err := geoip2.Open("data/GeoLite2-City.mmdb")
		if err != nil {
			fmt.Printf("Error initiating GeoIP2 database: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		defer db.Close()
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		ip := net.ParseIP(host)
		if ip == nil {
			fmt.Println("Invalid IP address")
			next.ServeHTTP(w, r)
			return
		}
		record, err := db.City(ip)
		if err != nil {
			fmt.Printf("Error getting GeoIP2 record: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		country := "Unknown"
		if name, ok := record.Country.Names["en"]; ok {
			country = name
		}

		city := "Unknown"
		if name, ok := record.City.Names["en"]; ok {
			city = name
		}
		referer := r.Referer()
		shortCode := vars["shortCode"]
		device := "desktop"
		ua := useragent.New(r.UserAgent())
		browser, _ := ua.Browser()
		os := ua.OS()
		if ua.Mobile() {
			device = "mobile"
		}
		if ua.Bot() {
			device = "bot"
		}
		node, err := snowflake.NewNode(1)
		if err != nil {
			fmt.Printf("failed to create snowflake node: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		id := node.Generate()
		clicksData := models.ClickAnalytics{
			ID:        id.String(),
			// Timestamp: time.Now().String(),
			ShortCode: shortCode,
			Referer:   referer,
			Ip:        hashIp(host),
			Country:   country,
			Os:        os,
			Browser:   browser,
			City:      city,
			Device:    device,
		}
		ctx := r.Context()
		clickKey := fmt.Sprintf("click:%s:%v", shortCode, clicksData.ID)
		hourStr := time.Now().Truncate(time.Hour).Format(time.RFC3339)
		pipeline := rdb.Pipeline()
		pipeline.HSet(ctx, clickKey,
			"ID", clicksData.ID,
			"Timestamp", time.Now().Format(time.RFC3339),
			"ShortCode", clicksData.ShortCode,
			"Referer", clicksData.Referer,
			"Ip", clicksData.Ip,
			"Country", clicksData.Country,
			"Os", clicksData.Os,
			"Browser", clicksData.Browser,
			"City", clicksData.City,
			"Device", clicksData.Device,
		)
		pipeline.Expire(ctx, clickKey, 24*time.Hour)
		pipeline.ZIncrBy(ctx, "clicks:"+shortCode+":by_hour", 1, hourStr)
		pipeline.ZIncrBy(ctx, "clicks:"+shortCode+":by_country", 1, country)
		pipeline.ZIncrBy(ctx, "clicks:"+shortCode+":by_city", 1, city)
		pipeline.ZIncrBy(ctx, "clicks:"+shortCode+":by_device", 1, device)
		pipeline.ZIncrBy(ctx, "clicks:"+shortCode+":by_referer", 1, referer)
		pipeline.XAdd(ctx, &redis.XAddArgs{
			Stream: "clicks:queue",
			Values: map[string]interface{}{
				"ID":        clicksData.ID,
				"Timestamp": time.Now().Format(time.RFC3339),
				"ShortCode": clicksData.ShortCode,
				"Referer":   clicksData.Referer,
				"Ip":        clicksData.Ip,
				"Country":   clicksData.Country,
				"Os":        clicksData.Os,
				"Browser":   clicksData.Browser,
				"City":      clicksData.City,
				"Device":    clicksData.Device,
			},
		})
		_, err = pipeline.Exec(ctx)
		if err != nil {
			fmt.Printf("Error executing pipeline: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
		// fmt.Printf("Click data: %+v\n", clicksData)
		// fmt.Printf("Click data stored in Redis: %+v\n", clickKey)
	})
}
