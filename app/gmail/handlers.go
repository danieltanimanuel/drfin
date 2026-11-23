package gmail

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"database/sql"
	"encoding/json"
	"errors"
	"log" 
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterHandler(r *mux.Router) {
	r.HandleFunc("/scan-email", h.scanEmail).Methods("POST")
}

func (h *Handler) scanEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Minute)
	defer cancel()

	var ownerID *int
	if v := r.URL.Query().Get("owner"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			http.Error(w, "invalid owner id", http.StatusBadRequest)
			return
		}
		ownerID = &id
	}

	results, err := h.service.SyncOwners(ctx, ownerID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			status = http.StatusNotFound
		}
		http.Error(w, fmt.Sprintf("email sync failed: %v", err), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	for _, res := range results {
		if !res.OK {
			w.WriteHeader(http.StatusInternalServerError) // any owner failing keeps a failing status
			break
		}
	}
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Error encoding sync results: %v", err)
	}
}
