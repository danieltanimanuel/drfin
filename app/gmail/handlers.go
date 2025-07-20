package gmail

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	auth *Auth
}

func NewHandler(auth *Auth) *Handler {
	return &Handler{
		auth: auth,
	}
}

func (h *Handler) RegisterHandler(r *mux.Router) {
	r.HandleFunc("/scan-email", h.scanEmail).Methods("POST")
}

func (h *Handler) scanEmail(w http.ResponseWriter, r *http.Request) {
	h.auth.authO2()
}
