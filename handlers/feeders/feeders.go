package feeders

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Rishabh-10/rss-agg/models"
	"github.com/go-chi/chi/v5"

	db "github.com/Rishabh-10/rss-agg/db/store"

	"github.com/google/uuid"
)

type Hander struct {
	queries *db.Queries
}

func New(q *db.Queries) *Hander {
	return &Hander{queries: q}
}

func (h *Hander) Create(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a feeder
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if len(data) == 0 {
		http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		return
	}

	var feeder models.Feeder
	err = json.Unmarshal(data, &feeder)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	feederCreated, err := h.queries.CreateFeeder(r.Context(), feeder.ToStoreModel())
	if err != nil {
		http.Error(w, "Failed to create feeder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(feederCreated)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func (h *Hander) GetFeeds(w http.ResponseWriter, r *http.Request) {
	// Implementation for Fetching feeds for a feeder
	var (
		limit  int
		offset int
	)

	id := chi.URLParam(r, "id")

	// validate id
	if !validateUUID(id) {
		http.Error(w, "Feeder ID is required", http.StatusBadRequest)
		return
	}

	// pagination
	limit, offset = getPaginationParams(r)

	var feeds []db.Feed

	feeds, err := h.queries.GetFeedsByFeederID(r.Context(), db.GetFeedsByFeederIDParams{
		FeederID: uuid.MustParse(id),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		http.Error(w, "Failed to fetch feeds: "+err.Error(), http.StatusInternalServerError)
		return
	}

	convertedFeeds := make([]models.FeedItem, 0, len(feeds))
	for _, feed := range feeds {
		convertedFeeds = append(convertedFeeds, models.FeedItem{
			ID:          feed.ID,
			FeederID:    feed.FeederID,
			Title:       feed.Title,
			Description: feed.Description,
		})
	}

	response := models.GetAll{
		Data: convertedFeeds,
		Meta: models.Meta{
			Limit:  limit,
			Offset: offset,
			Total:  len(convertedFeeds),
		},
	}

	data, _ := json.Marshal(response)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(data))
}

func validateUUID(id string) bool {
	if id == "" {
		return false
	}

	_, err := uuid.Parse(id)
	return err == nil
}

func getPaginationParams(r *http.Request) (limit int, offset int) {
	l, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || l == 0 {
		limit = 10
	}
	if l > 0 {
		limit = l
	}

	o, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || o == 0 {
		offset = 0
	}
	if o > 0 {
		offset = o
	}

	return limit, offset
}
