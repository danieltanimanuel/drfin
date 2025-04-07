package gmail

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterHandler(r *mux.Router) {
	r.HandleFunc("/scan-email", h.scanEmail).Methods("POST")
}

func (h *Handler) scanEmail(w http.ResponseWriter, r *http.Request) {
	authO2()
}
