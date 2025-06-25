package models

import (
	"github.com/google/uuid"

	db "github.com/Rishabh-10/rss-agg/db/store"
)

type Feed struct {
	Channel Channel `json:"channel" xml:"channel"`
}

type Channel struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" xml:"title"`
	Description string    `json:"description" xml:"description"`
}

func (f *Feed) ToStoreModel() db.CreateFeedParams {
	return db.CreateFeedParams{
		ID:          f.Channel.ID,
		Title:       f.Channel.Title,
		Description: f.Channel.Description,
	}
}
