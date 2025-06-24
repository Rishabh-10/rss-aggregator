package feeds

import (
	"net/http"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	// Implementation for fetching a feed
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Feed fetched"))
}
