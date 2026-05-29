package gmail

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	if err := h.service.syncMail(ctx); err != nil {
		http.Error(
			w,
			fmt.Sprintf("email sync failed: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sync completed"))
}
