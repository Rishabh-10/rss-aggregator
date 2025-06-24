package models

import (
	db "github.com/Rishabh-10/rss-agg/db/store"
	"github.com/google/uuid"
)

type Feeder struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// convert to store models
func (f *Feeder) ToStoreModel() db.CreateFeederParams {
	return db.CreateFeederParams{
		ID:   uuid.New(),
		Name: f.Name,
		Link: f.Link,
	}
}
