package dashboard

import (
	"net/http"
	"strconv"
	"drfin/app/owners"
)

const defaultOwnerID = 1

// ownerFromRequest returns ?owner=<id>, defaulting to owner 1. On a bad value
// it writes a 400 and returns ok=false.
func ownerFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	v := r.URL.Query().Get("owner")
	if v == "" {
		return defaultOwnerID, true
	}
	id, err := strconv.Atoi(v)
	if err != nil || id <= 0 {
		http.Error(w, "invalid owner id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

// ShellView is the data for dashboard.html.
type ShellView struct {
	Owners          []owners.Owner
	SelectedOwnerID int
	SelectedName    string
}