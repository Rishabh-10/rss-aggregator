package feedfetcher

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	db "github.com/Rishabh-10/rss-agg/db/store"
	"github.com/Rishabh-10/rss-agg/models"
	"github.com/google/uuid"
)

type FeedFetcher struct {
	DB             db.Queries
	WorkerPoolSize int
}

func New(db db.Queries, workerPoolSize int) *FeedFetcher {
	if workerPoolSize <= 0 {
		workerPoolSize = 5 // default value
	}

	return &FeedFetcher{
		DB:             db,
		WorkerPoolSize: workerPoolSize,
	}
}

func (f *FeedFetcher) GetFeeds() {
	slog.Info("Starting feed fetcher...")

	for {
		ctx := context.Background()

		feeders, err := f.DB.GetFeeders(ctx)
		if err != nil {
			slog.Error("failed to get feeders", "error", err)
		}

		if len(feeders) == 0 {
			slog.Info("no feeders found, retrying in 10 seconds")
			// Wait for 3 seconds before retrying
			<-time.After(3 * time.Second)
			continue
		}

		jobs := make(chan db.Feeder, len(feeders))
		var wg sync.WaitGroup

		for i := 0; i < f.WorkerPoolSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for feeder := range jobs {
					f.ProcessFeeders(ctx, feeder)
				}
			}()
		}

		for _, feeder := range feeders {
			jobs <- feeder
		}
		close(jobs)

		wg.Wait()

		slog.Info("Scrapped the DB, going into sleep for 1 min...")
		time.Sleep(1 * time.Minute)
	}
}

func (f *FeedFetcher) ProcessFeeders(ctx context.Context, feeder db.Feeder) {
	slog.Info("Fetching feed", "feeder", feeder.Name, "link", feeder.Link)

	res, err := http.Get(feeder.Link)
	if err != nil {
		slog.Error("failed to fetch feed", "feeder", feeder.Name, "error", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("non-200 status code received", "feeder", feeder.Name, "status_code", res.StatusCode)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("failed to read response body", "feeder", feeder.Name, "error", err)
		return
	}

	var feeds models.RSSFeeds

	_ = xml.Unmarshal(data, &feeds)

	slog.Info("Feed fetched successfully:" + fmt.Sprintf("%v", feeds))

	for _, feed := range feeds.Channel.Items {
		feed.ID = uuid.New()
		feed.FeederID = feeder.ID

		_, err := f.DB.CreateFeed(ctx, feed.ToStoreModel())
		if err != nil {
			slog.Error("failed to insert feed into database", "error", err, "feed_title", feed.Title)
			continue
		}

		slog.Info("Feed inserted successfully", "title", feed)
	}

	slog.Info("Hydrated all feeds for the Feeder: " + feeder.Link)
}
