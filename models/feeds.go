package models

import (
	db "github.com/Rishabh-10/rss-agg/db/store"
	"github.com/google/uuid"
)

type RSSFeeds struct {
	Channel struct {
		Title       string     `json:"title" xml:"title"`
		Description string     `json:"description" xml:"description"`
		Link        string     `json:"link" xml:"link"`
		Items       []FeedItem `json:"items" xml:"item"`
	} `json:"channel" xml:"channel"`
}

type FeedItem struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" xml:"title"`
	Description string    `json:"description" xml:"description"`
}

func (f *FeedItem) ToStoreModel() db.CreateFeedParams {
	return db.CreateFeedParams{
		ID:          f.ID,
		Title:       f.Title,
		Description: f.Description,
	}
}
