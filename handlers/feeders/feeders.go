package feeders

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Rishabh-10/rss-agg/models"

	db "github.com/Rishabh-10/rss-agg/db/store"
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
