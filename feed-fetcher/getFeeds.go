package feedfetcher

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Rishabh-10/rss-agg/db/store"
	"github.com/Rishabh-10/rss-agg/models"
	"github.com/google/uuid"
)

func GetFeeds(db db.Queries) {
	slog.Info("Starting feed fetcher...")

	for {
		ctx := context.Background()

		feeders, err := db.GetFeeders(ctx)
		if err != nil {
			slog.Error("failed to get feeders", "error", err)
		}

		if len(feeders) == 0 {
			slog.Info("no feeders found, retrying in 10 seconds")
			// Wait for 3 seconds before retrying
			<-time.After(3 * time.Second)
			continue
		}

		for _, feeder := range feeders {
			slog.Info("Fetching feed", "feeder", feeder.Name, "link", feeder.Link)

			res, err := http.Get(feeder.Link)
			if err != nil {
				slog.Error("failed to fetch feed", "feeder", feeder.Name, "error", err)
				continue
			}

			if res.StatusCode != http.StatusOK {
				slog.Error("non-200 status code received", "feeder", feeder.Name, "status_code", res.StatusCode)
				continue
			}

			data, err := io.ReadAll(res.Body)
			if err != nil {
				slog.Error("failed to read response body", "feeder", feeder.Name, "error", err)
				continue
			}

			var feeds models.RSSFeeds

			_ = xml.Unmarshal(data, &feeds)

			slog.Info("Feed fetched successfully:" + fmt.Sprintf("%v", feeds))

			for _, feed := range feeds.Channel.Items {
				feed.ID = uuid.New()
				feed.FeederID = feeder.ID

				_, err := db.CreateFeed(ctx, feed.ToStoreModel())
				if err != nil {
					slog.Error("failed to insert feed into database", "error", err, "feed_title", feed.Title)
					continue
				}

				slog.Info("Feed inserted successfully", "title", feed)
			}
		}
	}
}
