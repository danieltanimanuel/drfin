package accountgroups

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handler struct is responsible for handling HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates and returns a new Handler for transaction types
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterHandler registers the HTTP routes for AccountGroups
func (h *Handler) RegisterHandler(router *mux.Router) {
	// Route for GET request to retrieve all AccountGroups
	router.HandleFunc("/accountgroups", h.GetAccountGroups).Methods("GET")

	// Route for POST request to create a new accountgroup
	router.HandleFunc("/accountgroups", h.CreateAccountGroup).Methods("POST")

	// Route for GET request to retrieve a accountgroup by ID
	router.HandleFunc("/accountgroups/{id}", h.GetAccountGroupByID).Methods("GET")

	// Route for PUT request to update a accountgroup by ID
	router.HandleFunc("/accountgroups/{id}", h.UpdateAccountGroup).Methods("PUT")

	// Route for DELETE request to delete a accountgroup by ID
	router.HandleFunc("/accountgroups/{id}", h.DeleteAccountGroup).Methods("DELETE")
}

// GetTransactionAccountGroups handles GET requests to retrieve all transaction accountgroups
func (h *Handler) GetAccountGroups(w http.ResponseWriter, r *http.Request) {
	accountgroups, err := h.service.GetAllAccountGroups()
	if err != nil {
		http.Error(w, "Failed to retrieve AccountGroups", http.StatusInternalServerError)
		log.Printf("Error retrieving AccountGroups: %v", err)
		return
	}

	// Send the response wrapped in a success response format
	response := NewSuccessResponse(accountgroups)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// CreateAccountGroup handles POST requests to create a new accountgroup
func (h *Handler) CreateAccountGroup(w http.ResponseWriter, r *http.Request) {
	var accountgroup AccountGroup
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&accountgroup); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Call the service to create a new accountgroup
	createdAccountGroup, err := h.service.CreateAccountGroup(accountgroup)
	if err != nil {
		http.Error(w, "Failed to create accountgroup", http.StatusInternalServerError)
		log.Printf("Error creating accountgroup: %v", err)
		return
	}

	// Wrap the response in a standard "created" response structure
	response := NewCreatedResponse([]AccountGroup{createdAccountGroup})

	// Send JSON response with the created accountgroup
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// GetAccountGroupByID handles GET requests to retrieve a accountgroup by ID
func (h *Handler) GetAccountGroupByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid accountgroup ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to get the accountgroup by ID
	accountgroups, err := h.service.GetAccountGroupByID(idInt)
	if err != nil {
		http.Error(w, "AccountGroup not found", http.StatusNotFound)
		log.Printf("Error retrieving accountgroup by ID: %v", err)
		return
	}

	// Wrap the response in a standard success response structure
	response := NewSuccessResponse([]AccountGroup{*accountgroups})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateAccountGroup handles PUT requests to update an existing accountgroup by ID
func (h *Handler) UpdateAccountGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid accountgroup ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	var accountgroups AccountGroup
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&accountgroups); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Set the ID in the accountgroup (so the service knows which record to update)
	accountgroups.ID = idInt

	// Call the service to update the accountgroup
	updatedAccountGroup, err := h.service.UpdateAccountGroup(idInt, accountgroups)
	if err != nil {
		http.Error(w, "Failed to update accountgroup", http.StatusInternalServerError)
		log.Printf("Error updating accountgroup: %v", err)
		return
	}

	// Wrap the response in a success response format
	response := NewSuccessResponse([]AccountGroup{updatedAccountGroup})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteAccountGroup handles DELETE requests to delete a accountgroup by ID
func (h *Handler) DeleteAccountGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid accountgroup ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to delete the transaction accountgroup
	err = h.service.DeleteAccountGroup(idInt)
	if err != nil {
		http.Error(w, "Failed to delete accountgroup", http.StatusInternalServerError)
		log.Printf("Error deleting accountgroup: %v", err)
		return
	}

	// Send response indicating successful deletion
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
}
